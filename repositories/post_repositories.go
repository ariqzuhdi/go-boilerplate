package repositories

import (
	"errors"

	"github.com/cheeszy/journaling/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreatePost(db *gorm.DB, post *models.Post) error {
	return db.Create(post).Error
}

func FindPostByID(db *gorm.DB, id string) (*models.Post, error) {
	var post models.Post
	if err := db.First(&post, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func FindUserWithPostsByUsername(db *gorm.DB, username string) (*models.User, error) {
	var user models.User
	if err := db.Preload("Posts", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC")
	}).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdatePost(db *gorm.DB, post *models.Post) error {
	return db.Save(post).Error
}

func DeletePostByID(db *gorm.DB, id string) error {
	res := db.Where("id = ?", id).Delete(&models.Post{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("post not found or unauthorized")
	}
	return nil
}

func FindAllPosts(db *gorm.DB) ([]models.Post, error) {
	var posts []models.Post
	err := db.Find(&posts).Error
	return posts, err
}

func GetPostsByUserID(db *gorm.DB, userID uuid.UUID) ([]models.Post, error) {
	var posts []models.Post
	if err := db.Where("user_id = ?", userID).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}
