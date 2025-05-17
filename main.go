package main

import (
	"github.com/cheeszy/go-crud/controllers"
	"github.com/cheeszy/go-crud/initializers"
	"github.com/cheeszy/go-crud/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	router := gin.Default()
	// Users
	router.POST("/register", controllers.Register)
	router.GET("/users", controllers.Users)
	router.POST("/login", controllers.Login)

	// route need login
	authorized := router.Group("/")
	authorized.Use(middleware.RequireAuth)
	{
		authorized.POST("/posts", controllers.PostsCreate)

	}

	// Posts
	router.GET("/post/:id", controllers.PostsShow)
	router.GET("/posts", controllers.PostsIndex)
	router.PUT("/post/:id", controllers.PostsUpdate)
	router.DELETE("/post/:id", controllers.PostsDelete)
	router.GET("/monkeytype/", controllers.MonkeyAPI)
	router.NoRoute(controllers.NotFoundHandler)

	router.Run()
}
