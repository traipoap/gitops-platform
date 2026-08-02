package controllers

import (
	"exporter/models"
	"exporter/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Auth handlers
func HandleSignIn(db *gorm.DB, jwtService *services.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Extract Basic Auth credentials
		username, password, hasAuth := c.Request.BasicAuth()
		if !hasAuth || username == "" || password == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid credentials"})
			return
		}

		// 2. Retrieve user from SQLite database
		var user models.User
		// Query by username (email)
		if err := db.Where("username = ?", username).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			return
		}

		// 3. Verify password using bcrypt
		// Compare the input password with the stored hash in the database
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// 4. Generate JWT tokens using the retrieved user data
		// Convert user.ID (uint) to string for the JWT service
		userID := fmt.Sprintf("%d", user.ID)

		accessToken, err := jwtService.GenerateToken(userID, user.Username, user.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
			return
		}

		refreshToken, err := jwtService.GenerateRefreshToken(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
			return
		}

		// 5. Return tokens to client
		c.JSON(http.StatusOK, gin.H{
			"token":         accessToken,
			"refresh_token": refreshToken,
		})
	}
}
