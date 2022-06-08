package models

import (
	"gorm.io/gorm"
)

type Doc struct {
	gorm.Model
	FavoriteId   uint
	DocIndex     uint
	Url          string
	Summary 	 string
}

func (*Doc) TableName() string{
	return "docs"
}