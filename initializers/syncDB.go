package initializers

import (
	"log"
	"main/models"
)

func SyncDataBase() {
	errUser := DB.AutoMigrate(&models.User{})
	if errUser != nil {
		log.Fatal("Failed to AutoMigrate User!")
	}
}
