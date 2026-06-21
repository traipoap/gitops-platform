package services

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"exporter/models"

	"github.com/google/uuid"
)

type ExportService struct {
	quickwit    *QuickwitClient
	storagePath string   // local temp dir or S3 bucket
	pdpaFields  []string // fields to mask: ["source_ip", "user_id", ...]
}

func NewExportService(qw *QuickwitClient, storagePath string, pdpaFields []string) *ExportService {
	return &ExportService{
		quickwit:    qw,
		storagePath: storagePath,
		pdpaFields:  pdpaFields,
	}
}

// services/export_service.go - async version
func (s *ExportService) ExportAsync(ctx context.Context, params models.ExportParams) (string, error) {
	jobID := uuid.New().String()

	// Save job status to Redis/DB
	// s.jobStore.Set(jobID, JobStatus{Status: "processing", ...})

	// Run in background goroutine
	go func() {
		_, err := s.Export(context.Background(), params)
		// Update job status: ready/failed
		// s.jobStore.Update(jobID, result, err)
		// Optionally log the error
		if err != nil {
			fmt.Printf("Export job %s failed: %v\n", jobID, err)
		}
	}()

	return jobID, nil
}

// GET /export?async=true → returns { "job_id": "xxx", "status": "processing" }
// GET /export/status/:job_id → check progress
// GET /export/download/:job_id → download when ready

// escapeLucene escapes special characters in Lucene query syntax
// Reference: https://lucene.apache.org/core/2_9_4/queryparsersyntax.html#Escaping%20Special%20Characters
func escapeLucene(s string) string {
	chars := []rune{'\\', '+', '-', '!', '(', ')', ':', '^', '"', '{', '}', '[', ']', '}', '~', '*', '?', '|', '&', '/'}
	for _, ch := range chars {
		s = strings.ReplaceAll(s, string(ch), `\`+string(ch))
	}
	return s
}

// buildLuceneQuery constructs Lucene-style query from params
func buildLuceneQuery(p models.SearchParams) string {
	var parts []string

	// ✅ Timestamp range
	if p.FromTimestamp != nil && p.ToTimestamp != nil {
		parts = append(parts, fmt.Sprintf("timestamp:[%s TO %s]",
			*p.FromTimestamp, *p.ToTimestamp))
	}

	// ✅ Source IP with escaping
	if p.SourceIP != nil {
		parts = append(parts, fmt.Sprintf("source_ip:%s", escapeLucene(*p.SourceIP)))
	}

	// ✅ Custom filter with escaping
	if p.Message != nil {
		parts = append(parts, fmt.Sprintf("message:%s", escapeLucene(*p.Message)))
	}

	// ✅ Index timestamp (fix: ใช้ %d สำหรับ int64)
	if p.EndIndexTimestamp == nil {
		parts = append(parts, "index_timestamp:[* TO *]")
	} else {
		parts = append(parts, fmt.Sprintf("index_timestamp:[* TO %s]", *p.EndIndexTimestamp))
	}

	if len(parts) == 0 {
		return "*"
	}

	fmt.Println("query:", strings.Join(parts, " AND "))
	return strings.Join(parts, " AND ")
}

// Export executes the export job (sync mode for now)
func (s *ExportService) Export(ctx context.Context, params models.ExportParams) (*models.ExportResponse, error) {
	query := buildLuceneQuery(params.SearchParams) // reuse from search controller

	// Step 1: Fetch data from Quickwit (streaming approach)
	reader, recordCount, err := s.fetchAndTransform(ctx, query, params)
	// Only close if reader implements io.Closer
	if closer, ok := reader.(io.Closer); ok {
		defer closer.Close()
	}

	// Step 2: Prepare output file
	filename := params.Filename
	if filename == nil || *filename == "" {
		ts := time.Now().Format("20060102_150405")
		fn := fmt.Sprintf("export_%s.%s", ts, params.Format)
		filename = &fn
	}

	// Add extension if missing
	if !strings.HasSuffix(*filename, string(params.Format)) {
		*filename += "." + string(params.Format)
	}

	// Step 3: Write to storage (local file for now, can extend to S3)
	outputPath := filepath.Join(s.storagePath, *filename)

	// Handle gzip compression
	var writer io.Writer
	if params.Compress != nil && *params.Compress {
		outputPath += ".gz"
		file, err := os.Create(outputPath)
		if err != nil {
			return nil, fmt.Errorf("create file: %w", err)
		}
		gz := gzip.NewWriter(file)
		writer = gz
		defer func() { gz.Close(); file.Close() }()
	} else {
		file, err := os.Create(outputPath)
		if err != nil {
			return nil, fmt.Errorf("create file: %w", err)
		}
		writer = file
		defer file.Close()
	}

	// Step 4: Format and write data
	fileSize, err := s.writeData(writer, reader, params.Format)
	if err != nil {
		return nil, fmt.Errorf("write data: %w", err)
	}

	return &models.ExportResponse{
		DownloadURL: "/exports/" + filepath.Base(outputPath),
		Status:      "ready",
		RecordCount: recordCount,
		FileSize:    fileSize,
	}, nil
}

// fetchAndTransform streams data from Quickwit and applies transformations
func (s *ExportService) fetchAndTransform(ctx context.Context, query string, params models.ExportParams) (io.Reader, uint64, error) {
	// Use Quickwit client to search with pagination/streaming
	// For simplicity: fetch all in one batch (max_hits: 10000)
	// Production: implement cursor-based pagination or use Quickwit scroll API

	result, err := s.quickwit.Search(ctx, query, params.MaxHits, params.IndexID)
	if err != nil {
		return nil, 0, err
	}

	// Apply PDPA masking if requested
	if params.PDPACompliant != nil && *params.PDPACompliant {
		result.Hits = s.maskPDPAFields(result.Hits)
	}

	// Convert to reader for streaming write
	buf := &bytes.Buffer{}
	switch params.Format {
	case models.FormatNDJSON:
		for _, hit := range result.Hits {
			line, _ := json.Marshal(hit)
			buf.Write(line)
			buf.WriteByte('\n')
		}
	case models.FormatJSON:
		_ = json.NewEncoder(buf).Encode(result.Hits)
	case models.FormatCSV:
		if err := s.writeCSV(buf, result.Hits); err != nil {
			return nil, 0, err
		}
	}

	return buf, result.Total, nil
}

// maskPDPAFields redacts sensitive fields (simple implementation)
func (s *ExportService) maskPDPAFields(hits []map[string]interface{}) []map[string]interface{} {
	masked := make([]map[string]interface{}, len(hits))
	for i, hit := range hits {
		// Deep copy to avoid modifying original
		newHit := make(map[string]interface{}, len(hit))
		for k, v := range hit {
			if s.isPDPAField(k) {
				// Mask: show first 3 chars + "***"
				if str, ok := v.(string); ok && len(str) > 3 {
					newHit[k] = str[:3] + "***"
				} else {
					newHit[k] = "***"
				}
			} else {
				newHit[k] = v
			}
		}
		masked[i] = newHit
	}
	return masked
}

func (s *ExportService) isPDPAField(field string) bool {
	for _, f := range s.pdpaFields {
		if field == f || strings.HasPrefix(field, f+".") {
			return true
		}
	}
	return false
}

// writeCSV converts hits to CSV format
func (s *ExportService) writeCSV(w io.Writer, hits []map[string]interface{}) error {
	if len(hits) == 0 {
		return nil
	}

	// Collect all unique keys for headers
	headers := make(map[string]bool)
	for _, hit := range hits {
		for k := range hit {
			headers[k] = true
		}
	}
	headerList := make([]string, 0, len(headers))
	for k := range headers {
		headerList = append(headerList, k)
	}

	cw := csv.NewWriter(w)
	defer cw.Flush()

	// Write header
	if err := cw.Write(headerList); err != nil {
		return err
	}

	// Write rows
	for _, hit := range hits {
		row := make([]string, len(headerList))
		for i, col := range headerList {
			if val, ok := hit[col]; ok {
				row[i] = fmt.Sprintf("%v", val)
			}
		}
		if err := cw.Write(row); err != nil {
			return err
		}
	}
	return nil
}

// writeData writes formatted data to output and returns file size
func (s *ExportService) writeData(w io.Writer, r io.Reader, format models.ExportFormat) (int64, error) {
	// For streaming: copy directly
	// For in-memory buffer (current impl): just get size
	if buf, ok := r.(*bytes.Buffer); ok {
		n, err := buf.WriteTo(w)
		return n, err
	}
	return io.Copy(w, r)
}

// GetStoragePath returns the configured storage path for exports
func (s *ExportService) GetStoragePath() string {
	return s.storagePath
}
