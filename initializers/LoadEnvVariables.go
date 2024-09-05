package initializers

import (
	"github.com/joho/godotenv"
	"log"
)

func LoadEnVVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error while trying load .env")
	}
}
