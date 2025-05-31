package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Comment struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Content   string    `gorm:"type:text;not null"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	User      User
	PostID    uuid.UUID `gorm:"type:uuid;not null"`
	Post      Post
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
