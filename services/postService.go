package services

import (
	"time"

	"github.com/cheeszy/journaling/dto"
	"github.com/cheeszy/journaling/initializers"
	"github.com/cheeszy/journaling/models"
	"github.com/cheeszy/journaling/repositories"
	"github.com/google/uuid"
)

func CreatePost(req dto.CreatePostRequest, userID uuid.UUID) (*dto.PostResponse, error) {
	post := models.Post{
		Title:  req.Title,
		Body:   req.Body,
		UserID: userID,
	}

	if err := repositories.CreatePost(initializers.DB, &post); err != nil {
		return nil, err
	}

	return &dto.PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Body:      post.Body,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}, nil
}

func GetPostByID(id string) (*dto.PostResponse, error) {
	post, err := repositories.FindPostByID(initializers.DB, id)
	if err != nil {
		return nil, err
	}

	return &dto.PostResponse{
		ID:    post.ID,
		Title: post.Title,
		Body:  post.Body,
	}, nil
}

func GetPostsByUsername(username string) ([]dto.PostResponse, error) {
	user, err := repositories.FindUserWithPostsByUsername(initializers.DB, username)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.PostResponse, 0, len(user.Posts))
	for _, post := range user.Posts {
		responses = append(responses, dto.PostResponse{
			ID:        post.ID,
			Title:     post.Title,
			Body:      post.Body,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		})
	}

	return responses, nil
}

func UpdatePost(id string, req dto.UpdatePostRequest) (*dto.PostResponse, error) {
	post, err := repositories.FindPostByID(initializers.DB, id)
	if err != nil {
		return nil, err
	}

	post.Title = req.Title
	post.Body = req.Body
	post.UpdatedAt = time.Now()

	if err := repositories.UpdatePost(initializers.DB, post); err != nil {
		return nil, err
	}

	return &dto.PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Body:      post.Body,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}, nil
}

func DeletePost(id string) error {
	return repositories.DeletePostByID(initializers.DB, id)
}

func GetAllPosts() ([]models.Post, error) {
	return repositories.FindAllPosts(initializers.DB)
}
