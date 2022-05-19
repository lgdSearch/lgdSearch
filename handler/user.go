package handler

import (
	"lgdSearch/pkg/db"
	"lgdSearch/pkg/models"
	"gorm.io/gorm"
)

func QueryUser (id uint, username string) (*models.User, error){
	var user models.User
	var result *gorm.DB
	if(id != 0) {
		result = db.Engine.Select("id", "username", "nickname", "password").First(&user, id)
	}else {
		result = db.Engine.Select("id", "username", "nickname", "password").First(&user, "username = ?", username)
	}
	if result.Error != nil {
		return &user, result.Error
	}
	return &user, nil
}