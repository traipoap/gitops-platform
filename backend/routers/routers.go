package routers

import (
	"exporter/controllers"
	"exporter/middleware"
	"exporter/services"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
	if port := os.Getenv("PORT"); port != "" {
		return port
	}
	return defaultPort
}

// setupCORS configures CORS middleware
func setupCORS() cors.Config {
	return cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://192.168.1.107:3000", "http://localhost:4321"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
}

// setupAuthRoutes configures authentication routes
func setupAuthRoutes(r *gin.Engine, jwtService *services.JWTService) {
	// Public routes
	auth := r.Group("/api/auth")
	{
		auth.POST("/login", handleLogin(jwtService))
	}

	// Protected routes
	protected := r.Group("/api")
	protected.Use(middleware.JWTAuth())
	{
		protected.GET("/profile", controllers.HandleProfile)
		protected.GET("/admin", middleware.RoleAuth("admin"), controllers.HandleAdmin)
		protected.POST("/logout", controllers.HandleLogout)
	}
}

// setupAPIRoutes configures API endpoints
func setupAPIRoutes(r *gin.Engine) {
	api := r.Group("/api")
	api.Use(middleware.JWTAuth())
	{
		api.GET("/dashboard", controllers.HandleDashboard)
		api.GET("/data-inventory", controllers.HandleDataInventory)
		api.GET("/logs", controllers.HandleLogs)
		api.GET("/compliance", controllers.HandleCompliance)
		api.GET("/settings", controllers.HandleSettings)
		api.GET("/help", controllers.HandleHelp)
	}
}

// setupTaskRoutes configures task-related routes
func setupTaskRoutes(r *gin.Engine) {
	task := r.Group("/task")

	task.GET("/indices", controllers.HandleIndices)
	task.GET("/search", controllers.HandleSearch)
	task.GET("/export", controllers.HandleExport)
	task.GET("/exports", controllers.HandleExportsList)
	task.GET("/exports/:filename", controllers.HandleDownload)
	/*
	r.GET("/", HandleHome)
	r.LoadHTMLGlob("./pages/*.html")
	r.Static("/static", "./static")
	r.GET("/searcher", func(c *gin.Context) {
		c.HTML(http.StatusOK, "searcher.html", nil)
	})

	r.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	*/
}

// Auth handlers
func handleLogin(jwtService *services.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, password, hasAuth := c.Request.BasicAuth()

		// TODO: Move credentials to environment variables or use proper authentication
		if username == "admin@example.com" && password == "admin123" && hasAuth {
			token, err := jwtService.GenerateToken("0", "admin", "admin")
			refreshToken, err := jwtService.GenerateRefreshToken("0")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"token":         token,
				"refresh_token": refreshToken,
			})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}
}

// SetupRoutes configures all routes for the application
func SetupRoutes(r *gin.Engine) {
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
	setupAuthRoutes(r, jwtService)
	//setupAPIRoutes(r)
	setupTaskRoutes(r)

	// Start server
	r.Run(":" + port)
}
