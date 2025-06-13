package services

import (
	"net/http"

	"github.com/cheeszy/journaling/models"
	"github.com/gin-gonic/gin"
)

func CurrentUser(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	u := user.(models.User)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":       u.ID,
			"username": u.Username,
			"email":    u.Email,
		},
	})
}

func GetUserFromContext(c *gin.Context) (models.User, bool) {
	user, exists := c.Get("user")
	if !exists {
		return models.User{}, false
	}
	return user.(models.User), true
}
