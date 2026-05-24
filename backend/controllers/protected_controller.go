package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleProfile(c *gin.Context) {
	userID, _ := c.Get("userID")
	c.JSON(http.StatusOK, gin.H{"user_id": userID})
}

func HandleAdmin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Admin access granted"})
}

func HandleLogout(c *gin.Context) {
	// TODO: Implement token invalidation in production
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
