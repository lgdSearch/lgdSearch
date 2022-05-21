package models

import (
	"gorm.io/gorm"
)

type Favorite struct {
	gorm.Model
	UserId 		uint
	DocId 		uint
	Summary		string
}

func (*Favorite) TableName() string{
	return "favorites"
}