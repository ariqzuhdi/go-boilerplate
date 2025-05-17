package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title string `gorm:"not null"`
	Body  string `gorm:"type:text"`

	UserID uint // foreign key ke tabel users
	User   User // relasi dengan struct User
}
