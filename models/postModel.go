package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Post struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Title  string    `gorm:"not null" json:"title"`
	Body   string    `gorm:"type:text" json:"body"`
	UserID uuid.UUID `gorm:"type:uuid;not null" json:"userId"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
