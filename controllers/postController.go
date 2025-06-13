package controllers

import (
	"io"
	"net/http"
	"os"

	"github.com/cheeszy/journaling/dto"
	"github.com/cheeszy/journaling/initializers"
	"github.com/cheeszy/journaling/models"
	"github.com/gin-gonic/gin"
)

// HandleNotFound returns a 404 response for unregistered routes
func NotFoundHandler(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"error": "URL not found.",
	})
}

// PostsCreate handles creation of a new post
func PostsCreate(c *gin.Context) {
	var reqBody struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Build the post
	u := user.(models.User)
	post := models.Post{
		Title:  reqBody.Title,
		Body:   reqBody.Body,
		UserID: u.ID,
	}

	if err := initializers.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create post"})
		return
	}

	initializers.DB.Preload("User").First(&post, post.ID)

	postResponse := dto.PostResponse{
		ID:    post.ID,
		Title: post.Title,
		Body:  post.Body,
	}

	c.JSON(http.StatusOK, gin.H{
		"post": postResponse,
	})
}

// PostsShowById returns a single post by its ID
func PostsShowById(c *gin.Context) {
	id := c.Param("id")
	var post models.Post

	if err := initializers.DB.First(&post, "id = ?", id).Error; err != nil {
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

// PostsShowAllPosts returns all posts from a user by their username
func PostsShowAllPosts(c *gin.Context) {
	username := c.Param("username")
	var user models.User

	if err := initializers.DB.Preload("Posts").Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	postResponses := make([]dto.PostResponse, 0, len(user.Posts))
	for _, post := range user.Posts {
		postResponses = append(postResponses, dto.PostResponse{
			ID:    post.ID,
			Title: post.Title,
			Body:  post.Body,
		})
	}

	c.JSON(http.StatusOK, gin.H{"posts": postResponses})
}

// PostsUpdate modifies an existing post by ID
func PostsUpdate(c *gin.Context) {
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
	if err := initializers.DB.First(&post, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	initializers.DB.Model(&post).Updates(models.Post{
		Title: reqBody.Title,
		Body:  reqBody.Body,
	})

	postResponse := dto.PostResponse{
		ID:    post.ID,
		Title: reqBody.Title,
		Body:  reqBody.Body,
	}

	c.JSON(http.StatusOK, gin.H{"post": postResponse})
}

// PostsDelete removes a post by ID
func PostsDelete(c *gin.Context) {
	id := c.Param("id")
	var post models.Post

	result := initializers.DB.Where("id = ?", id).Delete(&post)
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

// PostsIndex returns all posts in the database (for testing only)
func PostsIndex(c *gin.Context) {
	var posts []models.Post
	initializers.DB.Find(&posts)

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
	})
}

// MonkeyAPI fetches data from MonkeyType's API
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
