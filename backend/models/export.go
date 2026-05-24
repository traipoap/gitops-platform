package models

type ExportFormat string

const (
	FormatJSON   ExportFormat = "json"
	FormatCSV    ExportFormat = "csv"
	FormatNDJSON ExportFormat = "ndjson" // newline-delimited JSON
)

type ExportParams struct {
	SearchParams               // embed search filters
	Format        ExportFormat `form:"format" binding:"oneof=json csv ndjson"`
	Filename      *string      `form:"filename"`
	Compress      *bool        `form:"compress" binding:"omitempty"` // gzip
	PDPACompliant *bool        `form:"pdpa" binding:"omitempty"`     // mask sensitive fields
}

type ExportResponse struct {
	JobID       string `json:"job_id,omitempty"`       // async export
	DownloadURL string `json:"download_url,omitempty"` // direct download
	Status      string `json:"status"`                 // "ready", "processing", "failed"
	RecordCount uint64 `json:"record_count"`
	FileSize    int64  `json:"file_size,omitempty"`
	Error       string `json:"error,omitempty"`
}
