package handler

import (
	"lgdSearch/pkg/db"
	"lgdSearch/pkg/models"
	"gorm.io/gorm"
	"time"
)

func QueryUser (id uint, username string) (*models.User, error){
	var user models.User
	var result *gorm.DB
	// 优先使用id查询
	if(id != 0) {
		result = db.Engine.First(&user, id)
	}else {
		result = db.Engine.First(&user, "username = ?", username)
	}
	if result.Error != nil {
		return &user, result.Error
	}
	return &user, nil
}

func CreateUser (username, password string) (uint, error) {
	user := models.User{
		Username: username,
		Password: password,
		Nickname: time.Now().String(),
	}
	result := db.Engine.Create(&user)
	return user.ID, result.Error
}

func UpdateUser (id uint, nickname string) error {
	result := db.Engine.Model(&models.User{Model: gorm.Model{ID: id}}).Updates(models.User{Nickname: nickname})
	return result.Error
}