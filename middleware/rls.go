package middleware

import (
	"net/http"

	"github.com/cheeszy/go-crud/initializers"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SetCurrentUserDB(db *gorm.DB, userID uuid.UUID) *gorm.DB {
	db.Exec("SET app.current_user_id = ?", userID.String())
	return db
}

func RequireRLS(c *gin.Context) {
	// Ambil userID dari context
	userID, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: User ID not found in context",
		})
		return
	}

	// Konversi userID ke uuid.UUID
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: Invalid user ID format",
		})
		return
	}

	// Set RLS untuk user yang sedang login
	db := SetCurrentUserDB(initializers.DB.Session(&gorm.Session{}), userUUID)

	c.Set("db", db)
	c.Next()
}
