package repositories

import (
	"github.com/cheeszy/journaling/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func FindUserByRecoveryKey(db *gorm.DB, recoveryKey string) (models.User, error) {
	var user models.User
	err := db.Where("recovery_key = ?", recoveryKey).First(&user).Error
	return user, err
}

func UpdateUsername(db *gorm.DB, userID uuid.UUID, newUsername string) error {
	return db.Model(&models.User{}).Where("id = ?", userID).Update("username", newUsername).Error
}

func UpdateEmail(db *gorm.DB, userID uuid.UUID, newEmail string) error {
	return db.Model(&models.User{}).Where("id = ?", userID).Update("email", newEmail).Error
}

func CreateUser(db *gorm.DB, user *models.User) error {
	return db.Create(user).Error
}

func FindUserByEmailOrUsername(db *gorm.DB, identifier string) (*models.User, error) {
	var user models.User
	err := db.Where("email = ? OR username = ?", identifier, identifier).First(&user).Error
	return &user, err
}

func FindUserByVerificationToken(db *gorm.DB, token string) (*models.User, error) {
	var user models.User
	err := db.Where("verification_token = ?", token).First(&user).Error
	return &user, err
}

func FindUserByEmail(db *gorm.DB, email string) (*models.User, error) {
	var user models.User
	err := db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func UpdateUser(db *gorm.DB, user *models.User) error {
	return db.Save(user).Error
}

func GetUserByEmail(db *gorm.DB, email string) (models.User, error) {
	var user models.User
	err := db.Where("email = ?", email).First(&user).Error
	return user, err
}
