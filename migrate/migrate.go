package main

import (
	"github.com/cheeszy/journaling/initializers"
	"github.com/cheeszy/journaling/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	initializers.DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})
}
