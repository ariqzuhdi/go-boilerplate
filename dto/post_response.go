package dto

import (
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	Username string `json:"username"`
}

type PostResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	// User  UserResponse `json:"user"`
}
