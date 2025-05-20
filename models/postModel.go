package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Post struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Title  string    `gorm:"not null"`
	Body   string    `gorm:"type:text"`
	UserID uuid.UUID `gorm:"type:uuid"` // foreign key
	User   User
}
