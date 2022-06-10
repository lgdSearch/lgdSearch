package handler

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"lgdSearch/pkg"
	"lgdSearch/pkg/db"
	"lgdSearch/pkg/models"
	colf "lgdSearch/pkg/utils/colf/doc"
)

func AppendFavorite(userId uint, name string) (uint, error) {
	favorite := models.Favorite{
		UserId: userId,
		Name:   name,
	}
	result := db.Engine.Create(&favorite)
	return favorite.ID, result.Error
}

func UpdateFavoriteName(userId, favId uint, name string) error {
	result := db.Engine.Model(&models.Favorite{Model: gorm.Model{ID: favId}, UserId: userId}).Updates(models.Favorite{Name: name})
	return result.Error
}

func DeleteFavorite(userId, favId uint) error {
	result := db.Engine.Select(clause.Associations).Where("user_id = ? AND id = ?", userId, favId).Delete(&models.Favorite{})
	return result.Error
}

func QueryFavorites(userId, limit, offset uint) ([]models.Favorite, int64, error) {
	favorites := make([]models.Favorite, 0, 10)
	err := db.Engine.Model(&models.User{Model: gorm.Model{ID: userId}}).Limit(int(limit)).Offset(int(offset)).Association("Favorites").Find(&favorites)
	if err != nil {
		return favorites, 0, err
	}
	total := db.Engine.Model(&models.User{Model: gorm.Model{ID: userId}}).Association("Favorites").Count()
	return favorites, total, err
}

func QueryFavorite(userId, favId uint, name string) (*models.Favorite, error) {
	var fav models.Favorite
	var result *gorm.DB
	// 优先使用id查询
	if favId != 0 {
		result = db.Engine.First(&fav, "id = ? AND user_id = ?", favId, userId)
	} else {
		result = db.Engine.First(&fav, "name = ? AND user_id = ?", name, userId)
	}
	if result.Error != nil {
		return &fav, result.Error
	}
	return &fav, nil
}

func AppendDoc(favId, docIndex uint) (*models.Doc, error) {
	buf := pkg.SearchEngine.GetDocById(uint32(docIndex))
	stDoc := colf.StorageIndexDoc{}
	err := stDoc.UnmarshalBinary(buf)
	if err != nil {
		return nil, err
	}
	doc := models.Doc{
		DocIndex: docIndex,
		Summary:  stDoc.Text,
		Url:      stDoc.Url,
	}
	if err = db.Engine.Model(&models.Favorite{Model: gorm.Model{ID: favId}}).Association("Docs").Append(&doc); err != nil {
		return nil, err
	}
	var docs models.Doc
	var result *gorm.DB
	result = db.Engine.First(&docs, "favorite_id = ? AND doc_index = ?", favId, docIndex)
	if result.Error != nil {
		return &docs, result.Error
	}
	return &docs, nil
}

func DeleteDoc(favId, docId uint) error {
	result := db.Engine.Where("favorite_id = ? AND id = ?", favId, docId).Delete(&models.Doc{})
	return result.Error
}

func QueryDoc(favId, docId, docIndex uint) (*models.Doc, error) {
	docs := models.Doc{}
	var result *gorm.DB
	if docId != 0 {
		result = db.Engine.First(&docs, "id = ? AND favorite_id = ?", docId, favId)
	} else {
		result = db.Engine.First(&docs, "doc_index = ? AND favorite_id = ?", docIndex, favId)
	}
	return &docs, result.Error
}

func QueryDocs(favId, limit, offset uint) ([]models.Doc, int64, error) {
	docs := make([]models.Doc, 0, 10)
	err := db.Engine.Model(&models.Favorite{Model: gorm.Model{ID: favId}}).Limit(int(limit)).Offset(int(offset)).Association("Docs").Find(&docs)
	if err != nil {
		return docs, 0, err
	}
	total := db.Engine.Model(&models.Favorite{Model: gorm.Model{ID: favId}}).Association("Docs").Count()
	return docs, total, err
}

func QueryFavoritesAndDocs(userId uint) ([]models.Favorite, error) {
	favorites := make([]models.Favorite, 0, 10)
	result := db.Engine.Preload("Docs").Find(&favorites, "user_id = ?", userId)
	return favorites, result.Error
}
