package handler

import (
	"lgdSearch/pkg/db"
	"lgdSearch/pkg/models"
	"gorm.io/gorm"
)

func AppendFavorite (userId, docId uint) error {
	favorite := models.Favorite{
		DocId: docId,
		// 生成summary，未完成
	}
	return db.Engine.Model(&models.User{Model: gorm.Model{ID: userId}}).Association("Favorites").Append(&favorite)
}

func DeleteFavorite (userId, docId uint) error {
	return db.Engine.Model(&models.User{Model: gorm.Model{ID: userId}}).Association("Favorites").Delete(&models.Favorite{DocId: docId})
}

func QueryFavorites (userId uint) ([]models.Favorite, error) {
	favorites := make([]models.Favorite, 0, 10)
	err := db.Engine.Model(&models.User{Model: gorm.Model{ID: userId}}).Association("Favorites").Find(&favorites)
	return favorites, err
}