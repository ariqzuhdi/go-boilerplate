package controllers

import (
	"io"
	"net/http"
	"os"

	"github.com/cheeszy/journaling/dto"
	"github.com/cheeszy/journaling/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NotFoundHandler(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"error": "URL not found.",
	})
}

func PostsCreate(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Bind input (sudah terenkripsi dari frontend)
	var reqBody struct {
		Title string `json:"title" binding:"required"` // terenkripsi
		Body  string `json:"body" binding:"required"`  // terenkripsi
	}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Auth
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	u := user.(models.User)

	// Simpan apa adanya (sudah terenkripsi)
	post := models.Post{
		Title:  reqBody.Title,
		Body:   reqBody.Body,
		UserID: u.ID,
	}

	if err := db.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"post": dto.PostResponse{
			ID:        post.ID,
			Title:     reqBody.Title, // masih terenkripsi
			Body:      reqBody.Body,  // masih terenkripsi
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		},
	})
}

func PostsShowById(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")
	var post models.Post

	if err := db.First(&post, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	postResponse := dto.PostResponse{
		ID:    post.ID,
		Title: post.Title,
		Body:  post.Body,
	}

	c.JSON(http.StatusOK, gin.H{"post": postResponse})
}

func PostsShowAllPosts(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	username := c.Param("username")

	// Auth check
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userInToken := currentUser.(models.User)

	if userInToken.Username != username {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	var user models.User
	if err := db.Preload("Posts", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC")
	}).Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	postResponses := make([]dto.PostResponse, 0, len(user.Posts))
	for _, post := range user.Posts {
		postResponses = append(postResponses, dto.PostResponse{
			ID:        post.ID,
			Title:     post.Title,
			Body:      post.Body,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": postResponses})
}

func PostsUpdate(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	var reqBody struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var post models.Post
	if err := db.First(&post, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	db.Model(&post).Updates(models.Post{
		Title: reqBody.Title,
		Body:  reqBody.Body,
	})

	postResponse := dto.PostResponse{
		ID:        post.ID,
		Title:     reqBody.Title,
		Body:      reqBody.Body,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{"post": postResponse})
}

func PostsDelete(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")
	var post models.Post

	result := db.Where("id = ?", id).Delete(&post)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found or unauthorized"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Deleted"})
}

func PostsIndex(c *gin.Context) {
	// NOTE: Ini hanya untuk testing global, jadi tidak pakai RLS
	db := c.MustGet("db").(*gorm.DB)
	var posts []models.Post
	db.Find(&posts)

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
	})
}

func MonkeyAPI(c *gin.Context) {
	apiKey := os.Getenv("MONKEYTYPE_API_KEY")

	req, _ := http.NewRequest("GET", "https://api.monkeytype.com/users/personalBests?mode=time", nil)
	req.Header.Add("Authorization", "ApeKey "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}

	c.Data(resp.StatusCode, "application/json", body)
}
