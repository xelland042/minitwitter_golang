package models

import "gorm.io/gorm"

type Tweet struct {
	gorm.Model
	Title    string `gorm:"column:title;not null"`
	Body     string `gorm:"column:body;not null"`
	File     string
	AuthorID uint `json:"author_id"`
	Author   User `gorm:"foreignKey:AuthorID"`
}
