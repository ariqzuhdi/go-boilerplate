package controllers

import (
	"github.com/cheeszy/go-crud/initializers"
	"github.com/cheeszy/go-crud/models"
	"github.com/gin-gonic/gin"
)

func PostsCreate(c *gin.Context) {
	// Get data off requests body
	var body struct {
		Body  string
		Title string
	}

	c.Bind(&body)

	post := models.Post{Title: body.Title, Body: body.Body}

	result := initializers.DB.Create(&post)

	if result.Error != nil {
		c.Status(400)
		return
	}

	// create a post

	// return it
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

func PostsShow(c *gin.Context) {
	// Get the id of url
	id := c.Param("id")

	var post []models.Post
	initializers.DB.First(&post, id)

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

func PostsDelete(c *gin.Context) {

	id := c.Param("id")
	var post []models.Post

	initializers.DB.Delete(&post, id)

	c.JSON(200, gin.H{
		"Message": "Deleted.",
	})
}
