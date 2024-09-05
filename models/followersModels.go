package models

import "gorm.io/gorm"

type FollowModel struct {
	gorm.Model
	Following    User
	FollowingID  uint `json:"following_id"`
	FollowedBy   User
	FollowedByID uint `json:"followed_id"`
}
