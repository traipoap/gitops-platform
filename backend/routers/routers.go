package routers

import (
	"exporter/controllers"
	"exporter/middleware"
	"exporter/services"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
// Same-origin proxying (npm run dev / dev:k3s) never triggers CORS; CORS_ORIGINS
// (comma-separated) only matters when a browser on another origin calls the
// backend directly, e.g. a local frontend testing the K3s backend:
//
//	CORS_ORIGINS="http://localhost:4321" go run main.go
func setupCORS() cors.Config {
	origins := []string{"http://localhost:4321", "http://127.0.0.1:4321", "https://frontend.example.com"}
	if v := strings.TrimSpace(os.Getenv("CORS_ORIGINS")); v != "" {
		origins = strings.Split(v, ",")
	}
	return cors.Config{
		AllowOrigins:     origins,
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
	r.GET("/metrics", HandleHome)
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
func SetupRoutes(r *gin.Engine, db *gorm.DB) {

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
