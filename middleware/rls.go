package middleware

import (
	"net/http"

	"github.com/cheeszy/journaling/initializers"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SetCurrentUserDB(db *gorm.DB, userID uuid.UUID) *gorm.DB {
	db.Exec("SET app.current_user_id = ?", userID.String())
	return db
}

func RequireRLS(c *gin.Context) {
	// grab userID from the context
	userID, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: User ID not found in context",
		})
		return
	}

	// convert userID to uuid.UUID
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: Invalid user ID format",
		})
		return
	}

	// set the current user in the database session
	db := SetCurrentUserDB(initializers.DB.Session(&gorm.Session{}), userUUID)

	c.Set("db", db)
	c.Next()
}
