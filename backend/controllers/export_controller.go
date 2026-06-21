package controllers

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
	exportDir := filepath.Join(os.TempDir(), "exports")
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		return err
	}

	// PDPA fields to mask (load from env/config)
	pdpaFields := []string{"source_ip", "user_id", "email", "phone"}

	exportService = services.NewExportService(qwClient, exportDir, pdpaFields)
	return nil
}

func HandleExport(c *gin.Context) {
	var params models.SearchParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters: " + err.Error()})
		return
	}

	query := buildLuceneQuery(params)
	fmt.Printf("Query: %s\n", query)
	result, err := quickwitClient.Search(c.Request.Context(), query, params.MaxHits, params.IndexID)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Search failed: " + err.Error()})
		return
	}
	// Make file
	csv_name := params.SourceIP
	if csv_name == nil {
		csv_name = new(string)
		*csv_name = "any"
	}
	now := time.Now()
	file, err := os.Create("exports/" + *csv_name + "_" + now.Format("20060102_150405") + ".csv")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	// Make Writer
	wtr := csv.NewWriter(file)
	defer wtr.Flush()

	data := [][]string{}
	total := result.Total
	processed := 0
	fmt.Println("Total:", total, "\nProcessed:", processed)

	for processed < int(total) {
		query := buildLuceneQuery(params)
		response, err := quickwitClient.Search(c.Request.Context(), query, params.MaxHits, params.IndexID)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "Search failed: " + err.Error()})
			return
		}
		if len(response.Hits) == 0 {
			break
		}

		for _, hit := range response.Hits {
			//source_ip := hit["source_ip"].(string)
			message := hit["message"].(string)
			data = append(data, []string{message})
			wtr.Write(data[len(data)-1])
			processed += 1
		}

		lastHit := response.Hits[len(response.Hits)-1]
		ts, ok := lastHit["index_timestamp"].(float64)
		strValue := strconv.FormatFloat(ts, 'f', 0, 64)
		params.EndIndexTimestamp = &strValue
		if !ok {
			fmt.Println("Error: index_timestamp is missing or not a number", ok)
			break
		}
		if ts == 0 {
			fmt.Println("index_timestamp is 0, stopping", ts)
			break
		}
	} // End Processed loop
	c.JSON(http.StatusOK, gin.H{"Total": total})
}

// Optional: Handler for downloading exported files
func HandleDownload(c *gin.Context) {
	// ─────────────────────────────────────────
	// STEP 1: รับค่าพารามิเตอร์จาก URL
	// ─────────────────────────────────────────
	// URL: /exports/my-log.json → filename = "my-log.json"
	filename := c.Param("filename")
	fmt.Println("downloading", filename)

	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "filename is required"})
		c.Abort() // หยุดการประมวลผลต่อ
		return
	}

	// ─────────────────────────────────────────
	// STEP 2: ป้องกัน Path Traversal Attack ⚠️
	// ─────────────────────────────────────────
	// Hacker อาจส่ง: ../../etc/passwd เพื่ออ่านไฟล์ระบบ!
	// filepath.Base() จะตัด path ออก เหลือแค่ชื่อไฟล์จริงๆ
	// "../../etc/passwd" → "passwd" ✅ ปลอดภัย
	cleanName := filepath.Base(filename)

	// ตรวจสอบเพิ่มเติม: ป้องกันชื่อไฟล์แปลกๆ
	if cleanName == "" || cleanName == "." || cleanName == ".." {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filename"})
		c.Abort()
		return
	}

	// ─────────────────────────────────────────
	// STEP 3: สร้างพาธเต็มของไฟล์ในระบบ
	// ─────────────────────────────────────────
	// สมมติ: storagePath = "/tmp/exports"
	// ผลลัพธ์: "/tmp/exports/my-log.json"
	filePath := filepath.Join("exports", cleanName)

	// ─────────────────────────────────────────
	// STEP 4: เช็คว่าไฟล์มีอยู่จริงและเป็นไฟล์ (ไม่ใช่โฟลเดอร์)
	// ─────────────────────────────────────────
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		// กรณี: ไฟล์ถูกลบไปแล้ว หรือชื่อผิด
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found or expired"})
		c.Abort()
		return
	}
	if err != nil {
		// กรณี: ไม่มีสิทธิ์อ่าน, disk error, อื่นๆ
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot access file"})
		c.Abort()
		return
	}
	if fileInfo.IsDir() {
		// ป้องกันการเข้าถึงโฟลเดอร์โดยบังเอิญ
		c.JSON(http.StatusBadRequest, gin.H{"error": "not a file"})
		c.Abort()
		return
	}

	// ─────────────────────────────────────────
	// STEP 5: ตั้งค่า HTTP Headers สำหรับดาวน์โหลด
	// ─────────────────────────────────────────

	// Header: บอก browser ว่านี่คือไฟล์สำหรับดาวน์โหลด
	// "attachment" = ให้ดาวน์โหลด, "inline" = ให้เปิดใน browser
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream") // binary data
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", cleanName))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	// Optional: บอกขนาดไฟล์ล่วงหน้า (browser แสดงความคืบหน้าได้)
	c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	// ─────────────────────────────────────────
	// STEP 6: ส่งไฟล์! 🚀
	// ─────────────────────────────────────────
	// Gin จะอ่านไฟล์และสตรีมไปยัง client โดยอัตโนมัติ
	// ใช้แรมน้อยแม้ไฟล์ใหญ่ เพราะไม่โหลดเข้าหน่วยความจำทั้งหมด
	c.File(filePath)

	// ไม่ต้อง return อะไรเพิ่ม เพราะ c.File() จัดการ response ให้แล้ว

}
