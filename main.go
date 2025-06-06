package main

import (
	"github.com/cheeszy/journaling/controllers"
	"github.com/cheeszy/journaling/initializers"
	"github.com/cheeszy/journaling/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	router := gin.Default()

	router.GET("/users", controllers.Users)
	// Users
	router.POST("/register", controllers.Register)
	router.POST("/login", controllers.Login)
	router.POST("/resend-verification", controllers.ResendVerificationEmail)

	// routes that need middleware
	authorized := router.Group("/")
	authorized.Use(middleware.RequireAuth)
	authorized.Use(middleware.RequireRLS)
	{
		authorized.POST("/posts", controllers.PostsCreate)
		authorized.GET("/post/:id", controllers.PostsShowById)
		authorized.GET("/:username", controllers.PostsShowAllPosts)
		authorized.PUT("/post/:id", controllers.PostsUpdate)
		authorized.DELETE("/post/:id", controllers.PostsDelete)
		authorized.POST("/logout", controllers.Logout)
	}

	// Posts
	router.GET("/monkeytype/", controllers.MonkeyAPI)
	router.NoRoute(controllers.NotFoundHandler)
	router.GET("/posts", controllers.PostsIndex)
	router.GET("/verify", controllers.VerifyEmail)

	router.Run()
}
