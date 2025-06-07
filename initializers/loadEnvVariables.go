package initializers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	if _, err := os.Stat(".env"); err == nil {
		// .env file exists, load it
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	} else {
		// .env file does not exist, skipping loading
		log.Println(".env file not found, skipping loading env file")
	}
}
