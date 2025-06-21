package services

import (
	"errors"
	"net/http"

	"github.com/cheeszy/journaling/dto"
	"github.com/cheeszy/journaling/initializers"
	"github.com/cheeszy/journaling/models"
	"github.com/cheeszy/journaling/repositories"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

func ResetPassword(input dto.ResetPasswordRequest) (*models.User, error) {
	db := initializers.DB
	user, err := repositories.FindUserByRecoveryKey(db, input.RecoveryKey)
	if err != nil {
		return nil, errors.New("invalid recovery key")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)

	if err := db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func ChangeUsername(userID uuid.UUID, newUsername string) error {
	db := initializers.DB
	return repositories.UpdateUsername(db, userID, newUsername)
}

func ChangeEmail(userID uuid.UUID, newEmail string) error {
	db := initializers.DB
	return repositories.UpdateEmail(db, userID, newEmail)
}
