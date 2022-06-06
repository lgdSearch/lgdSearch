package handler

import (
	"lgdSearch/pkg/db"
	"lgdSearch/pkg/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func AppendFavorite (userId uint, name string) (uint, error) {
	favorite := models.Favorite{
		UserId: userId,
		Name: name,
	}
	result := db.Engine.Create(&favorite)
	return favorite.ID, result.Error
}

func DeleteFavorite (favId uint) error {
	result := db.Engine.Select(clause.Associations).Delete(&models.Favorite{Model: gorm.Model{ID: favId}})
	return result.Error
}

func QueryFavorites (userId, limit, offset uint) ([]models.Favorite, error) {
	favorites := make([]models.Favorite, 0, 10)
	err := db.Engine.Model(&models.User{Model: gorm.Model{ID: userId}}).Limit(int(limit)).Offset(int(offset)).Association("Favorites").Find(&favorites)
	return favorites, err
}

func QueryFavorite (favId uint, name string) (*models.Favorite, error) {
	var fav models.Favorite
	var result *gorm.DB
	// 优先使用id查询
	if(favId != 0) {
		result = db.Engine.First(&fav, favId)
	}else {
		result = db.Engine.First(&fav, "name = ?", name)
	}
	if result.Error != nil {
		return &fav, result.Error
	}
	return &fav, nil
}

func AppendDoc (favId, docIndex uint) error {
	doc := models.Doc{
		DocIndex: docIndex,
		//summary生成 未完成
	}
	return db.Engine.Model(&models.Favorite{Model: gorm.Model{ID: favId}}).Association("Docs").Append(&doc)
}

func DeleteDoc (docId uint) error {
	result := db.Engine.Delete(&models.Doc{Model: gorm.Model{ID: docId}})
	return result.Error
}

func QueryDoc (docId uint, docIndex uint) (*models.Doc, error) {
	docs := models.Doc{}
	var result *gorm.DB
	if docId != 0 {
		result = db.Engine.First(&docs, docId)
	}else {
		result = db.Engine.First(&docs, "doc_index = ?", docIndex)
	}
	return &docs, result.Error
}

func QueryDocs (favId, limit, offset uint) ([]models.Doc, error) {
	docs := make([]models.Doc, 0, 10)
	err := db.Engine.Model(&models.Favorite{Model: gorm.Model{ID: favId}}).Limit(int(limit)).Offset(int(offset)).Association("Docs").Find(&docs)
	return docs, err
}