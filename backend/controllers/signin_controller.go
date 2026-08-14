package controllers

import (
	"exporter/models"
	"exporter/services"
	"fmt"
	"net/http"
	"strconv"

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
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// 4. Generate JWT tokens
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

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// HandleRefresh validates a refresh token and issues new access + refresh tokens.
func HandleRefresh(db *gorm.DB, jwtService *services.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req refreshRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing refresh_token"})
			return
		}

		// 1. Validate the refresh token
		userID, err := jwtService.ValidateRefreshToken(req.RefreshToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
			return
		}

		// 2. Look up user by ID
		var user models.User
		idUint, _ := strconv.ParseUint(userID, 10, 64)
		if err := db.First(&user, idUint).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		// 3. Issue new tokens (rotation: old refresh token is single-use by design)
		accessToken, err := jwtService.GenerateToken(userID, user.Username, user.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
			return
		}

		newRefreshToken, err := jwtService.GenerateRefreshToken(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token":         accessToken,
			"refresh_token": newRefreshToken,
		})
	}
}
