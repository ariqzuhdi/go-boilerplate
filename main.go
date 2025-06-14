package main

import (
	"time"

	"github.com/cheeszy/journaling/controllers"
	"github.com/cheeszy/journaling/initializers"
	"github.com/cheeszy/journaling/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Public endpoints
	api := router.Group("/api")
	api.POST("/register", controllers.Register)
	api.POST("/login", controllers.Login)
	api.POST("/resend-verification", controllers.ResendVerificationEmail)
	api.GET("/posts", controllers.PostsIndex) // all posts
	// api.GET("/posts/:id", controllers.PostsShowById) // by ID
	api.GET("/verify", controllers.VerifyEmail)
	api.GET("/monkeytype", controllers.MonkeyAPI)
	api.GET("/users", controllers.Users) // admin use

	// Protected endpoints
	authorized := router.Group("/api")
	authorized.Use(middleware.RequireAuth)
	authorized.Use(middleware.RequireRLS)
	{
		authorized.GET("/posts/user/:username", controllers.PostsShowAllPosts) // by username
		authorized.POST("/posts", controllers.PostsCreate)
		authorized.PUT("/posts/:id", controllers.PostsUpdate)
		authorized.DELETE("/posts/:id", controllers.PostsDelete)
		authorized.POST("/logout", controllers.Logout)
		authorized.GET("/", controllers.HomeHandler)
		authorized.GET("/user", controllers.GetCurrentUser)

	}

	router.NoRoute(controllers.NotFoundHandler)

	router.Run(":3000")
}
