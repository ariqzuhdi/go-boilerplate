package controllers

import (
	"net/http"

	"github.com/cheeszy/journaling/dto"
	"github.com/cheeszy/journaling/initializers"
	"github.com/cheeszy/journaling/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func ResetPasswordWithRecoveryKey(c *gin.Context) {
	var input dto.ResetPasswordRequest
	db := initializers.DB
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var user models.User
	if err := db.Where("recovery_key = ?", input.RecoveryKey).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid recovery key"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	// Optional: generate new session token
	c.JSON(http.StatusOK, gin.H{
		"message":  "Password reset successful",
		"username": user.Username,
	})

}

func ChangeUsername(c *gin.Context) {
	var req dto.ChangeUsernameRequest
	db := c.MustGet("db").(*gorm.DB)

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

	if err := models.UpdateUsername(db, userID, req.Username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update username", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Username updated successfully"})
}

func ChangeEmail(c *gin.Context) {
	var req dto.ChangeEmailRequest
	db := c.MustGet("db").(*gorm.DB)

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

	if err := models.UpdateEmail(db, userID, req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update email", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email updated successfully"})
}
