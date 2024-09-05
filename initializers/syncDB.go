package initializers

import (
	"log"
	"main/models"
)

func SyncDataBase() {
	errUser := DB.AutoMigrate(
		&models.User{},
		&models.Tweet{},
		&models.FollowModel{},
	)
	if errUser != nil {
		log.Fatal("Failed to AutoMigrate!")
	}
}
