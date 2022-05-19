package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID		  uint
	Username  string
	Nickname  string
	Password  string
	Favorites []uint
}

func (* User) TableName() string {
	return "User"
}