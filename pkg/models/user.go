package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string
	Nickname  string
	Password  string
	Favorites []Favorite
}

func (* User) TableName() string {
	return "users"
}