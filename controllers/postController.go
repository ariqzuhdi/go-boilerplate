package controllers

import (
	"io"
	"net/http"
	"os"

	"github.com/cheeszy/go-crud/initializers"
	"github.com/cheeszy/go-crud/models"
	"github.com/gin-gonic/gin"
)

func NotFoundHandler(c *gin.Context) {
	c.JSON(404, gin.H{
		"error": "URL not found.",
	})
}

func PostsCreate(c *gin.Context) {
	// Get data off requests body
	var body struct {
		Body  string
		Title string
	}

	c.Bind(&body)

	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(500, gin.H{"error": "Invalid user ID type"})
		return
	}

	post := models.Post{
		Title:  body.Title,
		Body:   body.Body,
		UserID: userID,
	}

	result := initializers.DB.Create(&post)
	if result.Error != nil {
		c.Status(400)
		return
	}

	initializers.DB.Preload("User").First(&post, post.ID)

	// return it
	c.JSON(200, gin.H{
		"post": post,
	})
}

func PostsShow(c *gin.Context) {

	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(500, gin.H{"error": "Invalid user ID type"})
		return
	}

	// Get the id of url
	id := c.Param("id")

	var post []models.Post
	request := initializers.DB.Where("id = ? AND user_id = ?", id, userID).First(&post, id)
	if request != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(200, gin.H{
		"post": post,
	})
}

func PostsUpdate(c *gin.Context) {
	// Get the id of url
	id := c.Param("id")

	// get the data off the req body

	var body struct {
		Body  string
		Title string
	}

	c.Bind(&body)

	// find the post were updateing
	var post []models.Post
	initializers.DB.First(&post, id)

	initializers.DB.Model(&post).Updates(models.Post{
		Title: body.Title, Body: body.Body,
	})

	// updating
	c.JSON(200, gin.H{
		"post": post,
	})
}

func PostsIndex(c *gin.Context) {
	//Get the posts
	var posts []models.Post
	initializers.DB.Find(&posts)

	//Respond with them
	c.JSON(200, gin.H{
		"post": posts,
	})
}

func PostsDelete(c *gin.Context) {

	id := c.Param("id")
	var post []models.Post

	initializers.DB.Delete(&post, id)

	c.JSON(200, gin.H{
		"Message": "Deleted.",
	})
}

func MonkeyAPI(c *gin.Context) {
	apiKey := os.Getenv("MONKEYTYPE_API_KEY")
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.monkeytype.com/users/personalBests?mode=time", nil)
	req.Header.Add("Authorization", "ApeKey "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to read response body"})
		return
	}

	// Forward response eksternal ke client
	c.Data(resp.StatusCode, "application/json", body)
}
