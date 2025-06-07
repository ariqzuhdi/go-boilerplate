package controllers

import "github.com/gin-gonic/gin"

func HomeHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Default home page",
		"status":  "YTTA",
	})
}
