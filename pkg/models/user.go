package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string `gorm:"size:30"`
	Nickname  string `gorm:"size:30"`
	Password  string `gorm:"size:20"`
	Favorites []Favorite
}

func (* User) TableName() string {
	return "users"
}