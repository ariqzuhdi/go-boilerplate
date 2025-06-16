package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/cheeszy/journaling/initializers"
	"github.com/cheeszy/journaling/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func RequireAuth(c *gin.Context) {
	var tokenString string

	// Coba ambil dari Authorization header
	authHeader := c.GetHeader("Authorization")

	if strings.HasPrefix(authHeader, "Bearer ") {
		tokenString = strings.TrimPrefix(authHeader, "Bearer ")
	} else {
		// Coba ambil dari cookie
		cookie, err := c.Cookie("token")
		if err == nil {
			tokenString = cookie
		}
	}

	if tokenString == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: Token not found",
		})
		return
	}

	// Lanjut validasi JWT seperti biasa
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: Invalid token",
		})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDStr, ok := claims["sub"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized: Invalid claims",
			})
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized: Invalid user ID format",
			})
			return
		}

		var user models.User
		if err := initializers.DB.First(&user, "id = ?", userID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized: User not found",
			})
			return
		}

		c.Set("user", user)
		c.Set("userID", userID)
		c.Next()
		return
	}

	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error": "Unauthorized: Invalid claims",
	})
}
