package models

import (
	"gorm.io/gorm"
)

type Doc struct {
	gorm.Model
	FavoriteId uint
	DocIndex   uint
	Url        string
	Summary    string
}

type DocId struct {
	DocId uint `json:"docId"`
}

func (*Doc) TableName() string {
	return "docs"
}
