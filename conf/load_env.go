package conf

import (
	"github.com/joho/godotenv"
	"log"
)

func init() {
	err := godotenv.Load("conf/.env")
	if err != nil {
		log.Println("No .env file found")
	}
}
