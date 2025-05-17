package main

import (
	"github.com/cheeszy/go-crud/controllers"
	"github.com/cheeszy/go-crud/initializers"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	router := gin.Default()
	router.POST("/posts", controllers.PostsCreate)
	router.GET("/post/:id", controllers.PostsShow)
	router.GET("/posts", controllers.PostsIndex)
	router.PUT("/post/:id", controllers.PostsUpdate)
	router.DELETE("/post/:id", controllers.PostsDelete)
	router.GET("/monkeytype/", controllers.MonkeyAPI)
	router.NoRoute(controllers.NotFoundHandler)
	router.Run()
}
