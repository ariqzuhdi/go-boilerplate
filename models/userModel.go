package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         uint   `gorm:"primarykey"`
	Username   string `gorm:"uniqueIndex;not null"`
	Email      string `gorm:"uniqueIndex;not null" json:"email" binding:"required, email"`
	Password   string `gorm:"not null" json:"password" binding:"required"`
	IsVerified bool   `gorm:"default:false"`

	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Posts []Post // satu user bisa punya banyak post
}
