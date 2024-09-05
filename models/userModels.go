package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName string `gorm:"column:username;unique"`
	Email    string `gorm:"unique"`
	Password string `gorm:"column:password;not null"`
	Bio      string
	Picture  string
	Tweets   []Tweet `gorm:"foreignKey:AuthorID"`
}
