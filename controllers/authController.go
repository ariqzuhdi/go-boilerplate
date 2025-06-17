package controllers

import (
	"net/http"

	"github.com/cheeszy/journaling/initializers"
	"github.com/cheeszy/journaling/models"
	"github.com/cheeszy/journaling/utils"
	"github.com/gin-gonic/gin"
)

type ResetPasswordRequest struct {
	Email                         string `json:"email" binding:"required,email"`
	RecoveryKey                   string `json:"recoveryKey" binding:"required"`
	NewPassword                   string `json:"newPassword" binding:"required,min=8"`
	EncryptedContentKeyByPassword string `json:"encryptedContentKeyByPassword" binding:"required"`
}

func ResetPasswordWithRecoveryKey(c *gin.Context) {
	var req ResetPasswordRequest
	user.EncryptedContentKeyByPassword = req.EncryptedContentKeyByPassword

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	db := initializers.DB

	// Ambil user berdasarkan email
	user, err := models.GetUserByEmail(db, req.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Cek recovery key
	if user.RecoveryKey != req.RecoveryKey {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid recovery key"})
		return
	}

	// // Ambil semua post milik user
	// posts, err := models.GetPostsByUserID(db, user.ID)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user posts"})
	// 	return
	// }

	// // Decrypt tiap post lalu encrypt ulang dengan password baru
	// for _, post := range posts {
	// 	decryptedTitle, err := utils.Decrypt(post.Title, req.RecoveryKey)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt title"})
	// 		return
	// 	}
	// 	decryptedBody, err := utils.Decrypt(post.Body, req.RecoveryKey)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt body"})
	// 		return
	// 	}

	// 	newTitle, err := utils.Encrypt(decryptedTitle, req.NewPassword)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt title"})
	// 		return
	// 	}
	// 	newBody, err := utils.Encrypt(decryptedBody, req.NewPassword)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt body"})
	// 		return
	// 	}

	// 	post.Title = newTitle
	// 	post.Body = newBody

	// 	if err := models.UpdatePost(db, &post); err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
	// 		return
	// 	}
	// }

	// Encrypt password baru
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = hashedPassword

	// Update user di database
	if err := models.UpdateUser(db, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successful"})
}
