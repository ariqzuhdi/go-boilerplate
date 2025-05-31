package dto

import "github.com/google/uuid"

type UserResponse struct {
	Username string `json:"username"`
}

type PostResponse struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
	Body  string    `json:"body"`
	// User  UserResponse `json:"user"`
}
