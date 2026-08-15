package controllers

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"exporter/config"
	"exporter/models"
	"exporter/services"

	"github.com/gin-gonic/gin"
)

var exportService *services.ExportService

func InitExportController() error {
	if err := config.Load(); err != nil {
		return err
	}

	qwClient := services.NewQuickwitClient(config.AppConfig.QuickwitURL)

	// Create export dir if not exists
	exportDir := filepath.Join(".", "exports")
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		return err
	}

	// PDPA fields to mask (load from env/config)
	pdpaFields := []string{"source_ip", "user_id", "email", "phone"}

	exportService = services.NewExportService(qwClient, exportDir, pdpaFields)
	return nil
}

const exportDirPath = "./exports"

// HandleExport runs the search against Quickwit and saves the results as a
// CSV file in ./exports, named "{sourceIP|any}_{YYYYMMDD_HHMMSS}.csv"
// (e.g. "192.168.1.1_20260401_213305.csv"). It returns the saved filename
// plus metadata so the frontend can display it and link to a download.
func HandleExport(c *gin.Context) {
	var params models.SearchParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters: " + err.Error()})
		return
	}

	// index_id is required (Quickwit routes the search by index).
	if params.IndexID == nil || *params.IndexID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "index_id is required"})
		return
	}

	// Default page size if the client didn't send one.
	maxHits := 10000
	if params.MaxHits != nil && *params.MaxHits > 0 {
		maxHits = *params.MaxHits
	}

	query := buildLuceneQuery(params)
	result, err := quickwitClient.Search(c.Request.Context(), query, &maxHits, params.IndexID)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Search failed: " + err.Error()})
		return
	}
	if len(result.Hits) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No logs matched the export query"})
		return
	}

	// ── Build CSV in memory ─────────────────────────────────────
	// Collect a stable header order (first-seen order across hits).
	headerIndex := make(map[string]int)
	var headers []string
	for _, hit := range result.Hits {
		for k := range hit {
			if _, ok := headerIndex[k]; !ok {
				headerIndex[k] = len(headers)
				headers = append(headers, k)
			}
		}
	}

	var buf bytes.Buffer
	wtr := csv.NewWriter(&buf)
	if err := wtr.Write(headers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write CSV header"})
		return
	}
	for _, hit := range result.Hits {
		row := make([]string, len(headers))
		for i, h := range headers {
			row[i] = cellString(hit[h])
		}
		if err := wtr.Write(row); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write CSV row"})
			return
		}
	}
	wtr.Flush()
	if err := wtr.Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to flush CSV"})
		return
	}

	// ── Save to ./exports ───────────────────────────────────────
	if err := os.MkdirAll(exportDirPath, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot create exports directory"})
		return
	}

	base := "any"
	if params.SourceIP != nil && *params.SourceIP != "" {
		base = sanitizeFileName(*params.SourceIP)
	}
	// Timestamp in Thai time (Asia/Bangkok, UTC+7, no DST) for the filename.
	bangkokLoc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		bangkokLoc = time.FixedZone("ICT", 7*60*60) // fallback: fixed UTC+7
	}
	now := time.Now().In(bangkokLoc)
	fileName := fmt.Sprintf("%s_%s.csv", base, now.Format("20060102_150405"))
	filePath := filepath.Join(exportDirPath, fileName)

	if err := os.WriteFile(filePath, buf.Bytes(), 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot write export file"})
		return
	}

	//c.JSON(http.StatusOK, gin.H{
	//	"filename": fileName,
	//	"path":     "exports/" + fileName,
	//	"total":    len(result.Hits),
	//	"size":     len(buf.Bytes()),
	//	"download": "/api/exports/" + fileName,
	//})

}

// cellString safely converts a Quickwit value to a CSV cell string.
func cellString(v interface{}) string {
	if v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(t)
	default:
		b, err := json.Marshal(t)
		if err != nil {
			return fmt.Sprintf("%v", v)
		}
		return string(b)
	}
}

var unsafeFileNameChars = regexp.MustCompile(`[^A-Za-z0-9._-]`)

// sanitizeFileName keeps only safe filename characters (so an IP like
// 192.168.1.1 stays intact, while spaces/slashes become underscores).
func sanitizeFileName(s string) string {
	return unsafeFileNameChars.ReplaceAllString(s, "_")
}

// Optional: Handler for downloading exported files
func HandleDownload(c *gin.Context) {
	filename := c.Param("filename")

	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "filename is required"})
		return
	}

	// Prevent path traversal: keep only the base name.
	cleanName := filepath.Base(filename)
	if cleanName == "" || cleanName == "." || cleanName == ".." {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filename"})
		return
	}

	filePath := filepath.Join("exports", cleanName)

	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found or expired"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot access file"})
		return
	}
	if fileInfo.IsDir() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "not a file"})
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", cleanName))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")
	c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	c.File(filePath)
}
