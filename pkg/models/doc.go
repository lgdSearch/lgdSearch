package models

import (
	"gorm.io/gorm"
)

type Doc struct {
	gorm.Model
	FavoriteId   uint
	DocIndex     uint
	Summary 	 string
}

func (*Doc) TableName() string{
	return "docs"
}