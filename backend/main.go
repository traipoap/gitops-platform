package main

import (
	"exporter/controllers"
	"exporter/models"
	"exporter/routers"
	"log"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Initialize Database (SQLite)
	db, err := gorm.Open(sqlite.Open("users.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Auto-migrate the User table
	db.AutoMigrate(&models.User{})

	var count int64
	db.Model(&models.User{}).Where("username = ?", "admin@example.com").Count(&count)
	if count == 0 {
		hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), 10)
		db.Create(&models.User{Username: "admin@example.com", PasswordHash: string(hash), Role: "admin"})
		log.Println("Seeded demo admin user")
	}

	if err := controllers.InitSearchController(); err != nil {
		log.Fatal("Failed to initialize: ", err)
	}
	if err := controllers.InitExportController(); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	routers.SetupRoutes(r, db)
	r.Run()
}
