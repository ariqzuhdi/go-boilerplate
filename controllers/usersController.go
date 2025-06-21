package controllers

import (
	"net/http"

	"github.com/cheeszy/journaling/dto"
	"github.com/cheeszy/journaling/initializers"
	"github.com/cheeszy/journaling/models"
	"github.com/cheeszy/journaling/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Register(c *gin.Context) {
	var input dto.RegisterRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recoveryKey, err := services.RegisterUser(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "User registered. Please check your email to verify your account.",
		"recoveryKey": recoveryKey,
	})
}

func Users(c *gin.Context) {
	var users []models.User
	initializers.DB.Find(&users)

	//Respond with them
	c.JSON(200, gin.H{
		"Users": users,
	})
}

func Login(c *gin.Context) {
	var input dto.LoginRequest
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	tokenString, expiresAt, err := services.LoginUser(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("token", tokenString, 3600*24, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"token":      tokenString,
		"expires_in": expiresAt,
	})
}
func Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", true, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
		return
	}

	if err := services.VerifyUserEmail(token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}
func GetCurrentUser(c *gin.Context) {
	user := c.MustGet("user")

	c.JSON(200, user)
}

func ResendVerificationEmail(c *gin.Context) {
	var input dto.ResendVerificationRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := services.ResendVerificationEmail(input.Email)
	if err != nil {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func ResetPasswordWithRecoveryKey(c *gin.Context) {
	var input dto.ResetPasswordRequest
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user, err := services.ResetPassword(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Password reset successful",
		"username": user.Username,
	})
}

func ChangeUsername(c *gin.Context) {
	var req dto.ChangeUsernameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}

	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID format"})
		return
	}

	if err := services.ChangeUsername(userID, req.Username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update username", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Username updated successfully"})
}

func ChangeEmail(c *gin.Context) {
	var req dto.ChangeEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}

	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID format"})
		return
	}

	if err := services.ChangeEmail(userID, req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update email", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email updated successfully"})
}
