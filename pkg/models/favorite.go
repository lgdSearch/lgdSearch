package models

import (
	"gorm.io/gorm"
)

type Favorite struct {
	gorm.Model
	UserId 		uint
	Name        string
	Docs 		[]Doc
}

func (*Favorite) TableName() string{
	return "favorites"
}