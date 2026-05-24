package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"exporter/config"
	"exporter/models"
	"exporter/services"

	"github.com/gin-gonic/gin"
)

var quickwitClient *services.QuickwitClient

func InitSearchController() error {
	if err := config.Load(); err != nil {
		return err
	}
	quickwitClient = services.NewQuickwitClient(config.AppConfig.QuickwitURL)
	return nil
}

func HandleSearch(c *gin.Context) {
	var params models.SearchParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters: " + err.Error()})
		return
	}
	query := buildLuceneQuery(params)
	result, err := quickwitClient.Search(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Search failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
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

	// ✅ Index timestamp (fix: ใช้ %s สำหรับ *string)
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

// escapeLucene escapes special characters in Lucene query syntax
// Reference: https://lucene.apache.org/core/2_9_4/queryparsersyntax.html#Escaping%20Special%20Characters
func escapeLucene(s string) string {
	chars := []rune{'\\', '+', '-', '!', '(', ')', ':', '^', '"', '{', '}', '[', ']', '}', '~', '*', '?', '|', '&', '/'}
	for _, ch := range chars {
		s = strings.ReplaceAll(s, string(ch), `\`+string(ch))
	}
	return s
}
