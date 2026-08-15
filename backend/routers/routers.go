package routers

import (
	"exporter/controllers"
	"exporter/middleware"
	"exporter/models"
	"exporter/services"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	defaultPort = "8080"
)

// HandleHome returns a welcome message with API usage information
func HandleHome(c *gin.Context) {
	port := getPort()
	c.String(http.StatusOK, "Dummy Quickwit API is running\n"+
		"Try: http://localhost:%s/api/v1/syslogs/search?query=8.8.8.8&max_hits=5", port)
}

// getPort returns the port from environment or default
func getPort() string {
	if port := os.Getenv("BACKEND_PORT"); port != "" {
		return port
	}
	return defaultPort
}

// setupCORS configures CORS middleware.
// Uses AllowOriginFunc to reflect the request's origin back, so the app works
// from localhost, 127.0.0.1, a LAN IP, a different dev port, or a tunnel —
// instead of only a hard-coded allow list (which caused silent
// "Failed to fetch" errors for any origin not in the list).
func setupCORS() cors.Config {
	return cors.Config{
		AllowOriginFunc: func(origin string) bool {
			// Allow any same-network dev origin. (Local dev tool — acceptable.)
			return origin != ""
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
}

// setupAuthRoutes configures authentication routes
func setupAuthRoutes(r *gin.Engine, jwtService *services.JWTService, db *gorm.DB) {
	// Public routes
	auth := r.Group("/api/auth")
	{
		auth.POST("/login", controllers.HandleSignIn(db, jwtService))
		auth.POST("/refresh", controllers.HandleRefresh(db, jwtService))
		//auth.POST("/register", controllers.HandleRegister(db))
	}

	// Protected routes
	api := r.Group("/api")
	api.Use(middleware.JWTAuth())
	{
		api.GET("/profile", controllers.HandleProfile)
		api.GET("/admin", middleware.RoleAuth("admin"), controllers.HandleAdmin)
		api.POST("/logout", controllers.HandleLogout)
	}
}

// setupAPIRoutes configures API-related routes
func setupAPIRoutes(r *gin.Engine, jwtService *services.JWTService) {
	r.GET("/", HandleHome)
	api := r.Group("/api")
	api.Use(middleware.JWTAuth())
	{
		api.GET("/indices", middleware.RoleAuth("admin"), controllers.HandleIndices)
		api.GET("/search", middleware.RoleAuth("admin"), controllers.HandleSearch)
		api.GET("/export", middleware.RoleAuth("admin"), controllers.HandleExport)
		api.GET("/exports", middleware.RoleAuth("admin"), controllers.HandleExportsList)
		api.DELETE("/exports/:filename", middleware.RoleAuth("admin"), controllers.HandleDeleteExport)
		api.GET("/exports/:filename", middleware.RoleAuth("admin"), controllers.HandleDownload)
	}
}

// SetupRoutes configures all routes for the application
func SetupRoutes(r *gin.Engine) {
	// Initialize Database (SQLite)
	db, err := gorm.Open(sqlite.Open("users.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Auto-migrate the User table
	db.AutoMigrate(&models.User{})

	if err := controllers.InitSearchController(); err != nil {
		log.Fatal("Failed to initialize: ", err)
	}
	if err := controllers.InitExportController(); err != nil {
		log.Fatal(err)
	}

	port := getPort()
	jwtService := services.NewJWTService()

	// Configure middleware
	r.Use(cors.New(setupCORS()))

	// Setup route groups
	setupAuthRoutes(r, jwtService, db)
	setupAPIRoutes(r, jwtService)

	// Start server
	r.Run(":" + port)
}
