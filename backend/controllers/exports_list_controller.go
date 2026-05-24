package controllers

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func HandleExportsList(c *gin.Context) {
	exportsDir := "./exports" // Change this to your actual exports directory if needed

	files, err := os.ReadDir(exportsDir)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to list exports"})
		return
	}

	var csvFiles []string
	for _, file := range files {
		if !file.IsDir() && (strings.HasSuffix(file.Name(), ".csv") || strings.HasSuffix(file.Name(), ".csv.gz") || strings.HasSuffix(file.Name(), ".zip")) {
			csvFiles = append(csvFiles, file.Name())
		}
	}

	c.JSON(200, gin.H{"exports": csvFiles})
}
