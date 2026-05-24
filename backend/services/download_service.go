package services

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// DownloadExportedCSV serves a CSV or CSV.GZ file for download by filename.
func DownloadExportedCSV(c *gin.Context, exportsDir string, filename string) {
	// Only allow .csv or .csv.gz files
	if !(filepath.Ext(filename) == ".csv" || filepath.Ext(filename) == ".gz") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
		return
	}

	fullPath := filepath.Join(exportsDir, filename)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.FileAttachment(fullPath, filename)
}
