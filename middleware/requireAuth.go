package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func RequireAuth(c *gin.Context) {
	// take the Authorization header from the request
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: Token missing",
		})
		return
	}

	// take token from the Authorization header
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Parse token dengan fungsi verifikasi secret
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// hmac method check
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		// secret key from environment variable
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: Invalid token",
		})
		return
	}

	// grab claims from the token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// grab user ID from claims
		// claims["sub"] is the subject, which we set to user ID when creating the token
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
		// save user ID to context for later use
		// this allows us to access user ID in handlers
		c.Set("userID", userID)

		// next middleware or handler
		c.Next()

	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: Invalid claims",
		})
	}
}
