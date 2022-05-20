package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string `gorm:"uniqueIndex"`
	Nickname  string
	Password  string
	Favorites []Favorite
}

func (* User) TableName() string {
	return "users"
}