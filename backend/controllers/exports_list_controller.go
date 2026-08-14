package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const exportsDir = "./exports"

type exportFile struct {
	Name    string `json:"name"`
	ModTime string `json:"time"`
	Size    uint64 `json:"size"`
	Hash    string `json:"hash"`
}

// HandleExportsList returns metadata for every export file (sorted newest-first).
func HandleExportsList(c *gin.Context) {
	files, err := os.ReadDir(exportsDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list exports"})
		return
	}

	var result []exportFile
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		name := f.Name()
		if !(strings.HasSuffix(name, ".csv") || strings.HasSuffix(name, ".csv.gz")) {
			continue
		}

		info, err := f.Info()
		if err != nil {
			continue
		}

		h, err := fileSHA256(filepath.Join(exportsDir, name))
		if err != nil {
			h = "error"
		}

		result = append(result, exportFile{
			Name:    name,
			ModTime: info.ModTime().Format(time.RFC3339),
			Size:    uint64(info.Size()),
			Hash:    fmt.Sprintf("SHA256:%s", h),
		})
	}

	// Sort newest first
	sort.Slice(result, func(i, j int) bool {
		return result[i].ModTime > result[j].ModTime
	})

	if result == nil {
		result = []exportFile{}
	}

	c.JSON(http.StatusOK, gin.H{"exports": result})
}

// HandleDeleteExport removes a single exported file.
func HandleDeleteExport(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "filename is required"})
		return
	}

	cleanName := filepath.Base(filename)
	if cleanName == "" || cleanName == "." || cleanName == ".." {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filename"})
		return
	}

	filePath := filepath.Join(exportsDir, cleanName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	if err := os.Remove(filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("deleted %s", cleanName)})
}

// fileSHA256 reads a file and returns its hex-encoded SHA-256 digest.
func fileSHA256(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
