package controllers

import (
	"net/http"

	"github.com/cheeszy/journaling/initializers"
	"github.com/cheeszy/journaling/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type ResetPasswordRequest struct {
	RecoveryKey string `json:"recoveryKey"`
	NewPassword string `json:"newPassword"`
}

func ResetPasswordWithRecoveryKey(c *gin.Context) {
	var input ResetPasswordRequest
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
