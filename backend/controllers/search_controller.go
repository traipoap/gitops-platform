package controllers

import (
	"encoding/json"
	"fmt"
	"io"
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

func HandleIndices(c *gin.Context) {
	// Reuse the shared keep-alive client instead of allocating a new
	// http.Client (and a new TLS handshake) per request.
	resp, err := quickwitClient.Get(c.Request.Context(), "/api/v1/indexes")
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Cannot connect to Quickwit: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	// อ่านข้อมูลดิบ (Raw Body) ที่ส่งกลับมาจาก Quickwit
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read Quickwit response: " + err.Error()})
		return
	}

	// ตรวจสอบสถานะการตอบกลับจาก Quickwit
	if resp.StatusCode != http.StatusOK {
		c.Data(resp.StatusCode, "application/json", bodyBytes)
		return
	}

	// แปลงข้อมูล JSON ดิบเพื่อส่งต่อให้ฝั่ง Client ของคุณทันที
	var quickwitResponse []interface{} // Quickwit คืนค่ากลับมาเป็น Array ของ Index Object
	if err := json.Unmarshal(bodyBytes, &quickwitResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse JSON: " + err.Error()})
		return
	}

	// ส่งผลลัพธ์กลับในรูปแบบ JSON สวยงาม
	c.JSON(http.StatusOK, quickwitResponse)
}

func HandleSearch(c *gin.Context) {
	var params models.SearchParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters: " + err.Error()})
		return
	}

	query := buildLuceneQuery(params)
	result, err := quickwitClient.Search(c.Request.Context(), query, params.MaxHits, params.IndexID)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Search failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// luceneReplacer escapes Lucene special characters in a single pass.
// (The old per-character strings.ReplaceAll loop copied the string ~18 times.)
var luceneReplacer = strings.NewReplacer(
	`\`, `\\`,
	`+`, `\+`, `-`, `\-`, `!`, `\!`, `(`, `\(`, `)`, `\)`,
	`:`, `\:`, `^`, `\^`, `"`, `\"`,
	`{`, `\{`, `}`, `\}`, `[`, `\[`, `]`, `\]`,
	`~`, `\~`, `*`, `\*`, `?`, `\?`, `|`, `\|`, `&`, `\&`, `/`, `\/`,
)

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

	return strings.Join(parts, " AND ")
}

// escapeLucene escapes special characters in Lucene query syntax
// Reference: https://lucene.apache.org/core/2_9_4/queryparsersyntax.html#Escaping%20Special%20Characters
func escapeLucene(s string) string {
	return luceneReplacer.Replace(s)
}
