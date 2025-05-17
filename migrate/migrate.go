package main

import (
	"github.com/cheeszy/go-crud/initializers"
	"github.com/cheeszy/go-crud/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	initializers.DB.AutoMigrate(&models.User{}, &models.Post{})
}
