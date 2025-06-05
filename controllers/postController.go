package controllers

import (
	"io"
	"net/http"
	"os"

	"github.com/cheeszy/go-crud/dto"
	"github.com/cheeszy/go-crud/initializers"
	"github.com/cheeszy/go-crud/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func NotFoundHandler(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"error": "URL not found.",
	})
}

func PostsCreate(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	// Get data off requests body
	var body struct {
		Title string
		Body  string
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userIDAny, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID := userIDAny.(uuid.UUID)

	post := models.Post{
		Title:  body.Title,
		Body:   body.Body,
		UserID: userID,
	}

	if err := db.Create(&post).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create post"})
		return
	}

	db.Preload("User").First(&post, post.ID)

	postResponse := dto.PostResponse{
		ID:    post.ID,
		Title: post.Title,
		Body:  post.Body,
	}

	// return it
	c.JSON(200, gin.H{
		"post": postResponse,
	})
}

func PostsShowById(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	var post models.Post

	result := db.First(&post, "id = ?", id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	postResponse := dto.PostResponse{
		ID:    post.ID,
		Title: post.Title,
		Body:  post.Body, // Initialize with zero value
	}

	c.JSON(200, gin.H{
		"post": postResponse,
	})
}

func PostsShowAllPosts(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)
	username := c.Param("username")

	// take user from username param
	var user models.User

	if err := db.Preload("Posts").Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var postResponses []dto.PostResponse
	for _, post := range user.Posts {
		postResponses = append(postResponses, dto.PostResponse{
			ID:    post.ID,
			Title: post.Title,
			Body:  post.Body,
			// User:  dto.UserResponse{Username: post.User.Username},
		})
	}

	c.JSON(200, gin.H{
		"posts": postResponses,
	})
}

func PostsUpdate(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	// Get the id of url
	id := c.Param("id")

	// get the data off the req body

	var body struct {
		Body  string
		Title string
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// find the post were updating
	var post models.Post
	result := db.Where("id = ?", id).First(&post)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	db.Where("id = ?", id).Model(&post).Updates(models.Post{
		Title: body.Title, Body: body.Body,
	})

	postResponse := dto.PostResponse{
		ID:    post.ID,
		Title: post.Title,
		Body:  post.Body,
	}

	c.JSON(200, gin.H{
		"post": postResponse,
	})
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

// BE CAREFUL: THIS IS NOT PROTECTED
// It is just for testing purposes
func PostsIndex(c *gin.Context) {
	//Get the posts
	var posts []models.Post
	initializers.DB.Find(&posts)

	//Respond with them
	c.JSON(http.StatusAccepted, gin.H{
		"post": posts,
	})
}

func MonkeyAPI(c *gin.Context) {
	apiKey := os.Getenv("MONKEYTYPE_API_KEY")
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.monkeytype.com/users/personalBests?mode=time", nil)
	req.Header.Add("Authorization", "ApeKey "+apiKey)

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

	// forward the response
	c.Data(resp.StatusCode, "application/json", body)
}
