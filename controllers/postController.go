package controllers

import (
	"io"
	"net/http"
	"os"

	"github.com/cheeszy/journaling/dto"
	"github.com/cheeszy/journaling/models"
	"github.com/cheeszy/journaling/services"
	"github.com/gin-gonic/gin"
)

func NotFoundHandler(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"error": "URL not found.",
	})
}

func PostsCreate(c *gin.Context) {
	var req dto.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	u := user.(models.User)

	post, err := services.CreatePost(req, u.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"post": post})
}

func PostsShowById(c *gin.Context) {
	id := c.Param("id")
	post, err := services.GetPostByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"post": post})
}

func PostsShowAllPosts(c *gin.Context) {
	username := c.Param("username")
	userInToken := c.MustGet("user").(models.User)

	if userInToken.Username != username {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	posts, err := services.GetPostsByUsername(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": posts})
}

func PostsUpdate(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	post, err := services.UpdatePost(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"post": post})
}

func PostsDelete(c *gin.Context) {
	id := c.Param("id")

	if err := services.DeletePost(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Deleted"})
}

func PostsIndex(c *gin.Context) {
	posts, err := services.GetAllPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"posts": posts})
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
