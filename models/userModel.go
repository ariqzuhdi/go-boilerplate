package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID                     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Username               string    `gorm:"uniqueIndex;not null" json:"username"`
	Email                  string    `gorm:"uniqueIndex;not null" json:"email" binding:"required,email"`
	Password               string    `gorm:"not null" json:"password" binding:"required"`
	ResendCount            int       `gorm:"default:0" json:"resendCount,omitempty"`
	VerificationToken      string    `gorm:"default:null" json:"verificationToken,omitempty"`
	LastVerificationSentAt time.Time `gorm:"default:null" json:"lastVerificationSentAt,omitempty"`
	VerificationExpiresAt  time.Time `gorm:"default:null" json:"verificationExpiresAt,omitempty"`
	IsVerified             bool      `gorm:"default:false" json:"isVerified"`

	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Posts []Post `gorm:"foreignKey:UserID" json:"posts,omitempty"`
}
