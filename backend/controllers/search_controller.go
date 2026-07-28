package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

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
	// 1. สร้าง HTTP Client พร้อมตั้งค่า Timeout ป้องกันระบบค้าง
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// 2. สร้างคำสั่งยิงไปที่ Quickwit REST API (/api/v1/indexes)
	quickwitEndpoint := strings.TrimSuffix(config.AppConfig.QuickwitURL, "/") + "/api/v1/indexes"
	req, err := http.NewRequestWithContext(c.Request.Context(), "GET", quickwitEndpoint, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request: " + err.Error()})
		return
	}

	// 3. เริ่มส่ง Request ไปยัง Quickwit
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Cannot connect to Quickwit: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	// 4. อ่านข้อมูลดิบ (Raw Body) ที่ส่งกลับมาจาก Quickwit
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read Quickwit response: " + err.Error()})
		return
	}

	// 5. ตรวจสอบสถานะการตอบกลับจาก Quickwit
	if resp.StatusCode != http.StatusOK {
		c.Data(resp.StatusCode, "application/json", bodyBytes)
		return
	}

	// 6. แปลงข้อมูล JSON ดิบเพื่อส่งต่อให้ฝั่ง Client ของคุณทันที
	var quickwitResponse []interface{} // Quickwit คืนค่ากลับมาเป็น Array ของ Index Object
	if err := json.Unmarshal(bodyBytes, &quickwitResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse JSON: " + err.Error()})
		return
	}

	// 7. ส่งผลลัพธ์กลับในรูปแบบ JSON สวยงาม
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

// buildLuceneQuery constructs Lucene-style query from params
func buildLuceneQuery(p models.SearchParams) string {
	var parts []string

	//fmt.Printf("first query: %+v\n", p)

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

	//fmt.Println("last query:", strings.Join(parts, " AND "))
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
