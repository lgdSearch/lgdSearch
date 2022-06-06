package handler_test

import (
	"errors"
	"lgdSearch/handler"
	"lgdSearch/pkg/db"
	"lgdSearch/pkg/models"
	"testing"
	"strings"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	viper.SetConfigFile("test.env")
	if err := viper.ReadInConfig(); err != nil {
		panic("env initialization failed")
	}
	var err error
	db.Engine, err = gorm.Open(mysql.Open(viper.GetString("mysql_dsn")))
	if err != nil {
		panic("failed to connect database")
	}
	db.Engine.AutoMigrate(&models.User{})
	db.Engine.AutoMigrate(&models.Favorite{})
	db.Engine.AutoMigrate(&models.Doc{})
	m.Run()
	db.Engine.Where("1 = 1").Delete(&models.User{})
	db.Engine.Where("1 = 1").Delete(&models.Favorite{})
	db.Engine.Where("1 = 1").Delete(&models.Doc{})
}

func newUserId(name string) uint {
	user := models.User{
		Username: name + "_test",
		Password: "test",
	}
	db.Engine.Create(&user)
	return user.ID
}

func newFavoriteId(userId uint, name string) uint {
	favorite := models.Favorite{
		UserId: userId,
		Name: name,
	}
	db.Engine.Create(&favorite)
	return favorite.ID
}

func newDocID(favId uint, docIndex uint) uint {
	doc := models.Doc{
		FavoriteId: favId,
		DocIndex: docIndex,
	}
	db.Engine.Create(&doc)
	return doc.ID
}

func TestCreateUser(t *testing.T) {
	_, err := handler.CreateUser("test", "123456")
	if err != nil {
		t.Error(err.Error())
	}
}

func TestQueryUser(t *testing.T) {
	id := newUserId("QueryUser")
	_, err := handler.QueryUser(id, "")
	if err != nil {
		t.Error(err.Error())
	}
	_, err = handler.QueryUser(0, "QueryUser_test")
	if err != nil {
		t.Error(err.Error())
	}
}

func TestUpdateUserNickname(t *testing.T) {
	id := newUserId("UpdateUserNickname")
	err := handler.UpdateUserNickname(id, "test")
	if err != nil {
		t.Error(err.Error())
	}
	rUser := models.User{}
	result := db.Engine.First(&rUser, id)
	if result.Error != nil {
		t.Error(result.Error.Error())
	}
	if rUser.Nickname != "test" {
		t.Error("nickname incorrect")
	}
}

func TestDeleteUser(t *testing.T) {
	id := newUserId("DeleteUser")
	err := handler.DeleteUser(id)
	if err != nil {
		t.Error(err.Error())
	}
	user := models.User{}
	result := db.Engine.First(&user, id)
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		if result.Error == nil {
			t.Error("entry not deleted")
		} else {
			t.Error(result.Error.Error())
		}
	}
}

func TestAppendFavorite(t *testing.T) {
	id := newUserId("AppendFavorite")
	_, err := handler.AppendFavorite(id, "TestAppendFavorite")
	if err != nil {
		t.Error(err.Error())
	}
}

func TestDeleteFavorite(t *testing.T) {
	userId := newUserId("DeleteFavorite")
	favId := newFavoriteId(userId, "DeleteFavorite")
	newFavoriteId(userId, "DeleteFavorite")
	newFavoriteId(userId, "DeleteFavorite")
	err := handler.DeleteFavorite(favId)
	if err != nil {
		t.Error(err.Error())
	}
	num := db.Engine.Model(&models.User{Model: gorm.Model{ID: userId}}).Association("Favorites").Count()
	if num != 2 {
		t.Error("incorrect number of entries")
	}
}

func TestQueryFavorites(t *testing.T) {
	id := newUserId("QueryFavorites")
	newFavoriteId(id, "QueryFavorites")
	newFavoriteId(id, "QueryFavorites")
	newFavoriteId(id, "QueryFavorites")
	result, err := handler.QueryFavorites(id, 10, 0)
	if err != nil {
		t.Error(err.Error())
	}
	if len(result) != 3 {
		t.Error("incorrect number of entries")
	}
	result, err = handler.QueryFavorites(id, 10, 1)
	if err != nil {
		t.Error(err.Error())
	}
	if len(result) != 2 {
		t.Error("incorrect number of entries")
	}
	result, err = handler.QueryFavorites(id, 1, 0)
	if err != nil {
		t.Error(err.Error())
	}
	if len(result) != 1 {
		t.Error("incorrect number of entries")
	}
}

func TestQueryFavorite(t *testing.T) {
	userId := newUserId("QueryFavorite")
	favId := newFavoriteId(userId, "TestQueryFavorite")
	fav, err := handler.QueryFavorite(favId, "")
	if err != nil {
		t.Error(err.Error())
	}
	if strings.Compare(fav.Name, "TestQueryFavorite") != 0 {
		t.Error("name incorrect")
	}
	fav, err = handler.QueryFavorite(0, "TestQueryFavorite")
	if err != nil {
		t.Error(err.Error())
	}
	if fav.ID != favId {
		t.Error("id incorrect")
	}
}

func TestAppendDoc(t *testing.T) {
	userId := newUserId("TestAppendDoc")
	favId := newFavoriteId(userId, "TestAppendDoc")
	err := handler.AppendDoc(favId, 1)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestDeleteDoc(t *testing.T) {
	userId := newUserId("TestDeleteDoc")
	favId := newFavoriteId(userId, "TestDeleteDoc")
	docId := newDocID(favId, 3)
	err := handler.DeleteDoc(docId)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestQueryDoc(t *testing.T) {
	userId := newUserId("TestQueryDoc")
	favId := newFavoriteId(userId, "TestQueryDoc")
	docId := newDocID(favId, 3)
	doc, err := handler.QueryDoc(docId, 0)
	if err != nil {
		t.Error(err.Error())
	}
	if doc.DocIndex != 3 {
		t.Error("doc_index incorrect")
	}
	doc, err = handler.QueryDoc(0, 3)
	if err != nil {
		t.Error(err.Error())
	}
	if doc.ID != docId {
		t.Error("doc_id incorrect")
	}
}

func TestQueryDocs( t *testing.T) {
	userId := newUserId("TestQueryDoc")
	favId := newFavoriteId(userId, "TestQueryDoc")
	newDocID(favId, 1)
	newDocID(favId, 2)
	newDocID(favId, 3)
	docs, err := handler.QueryDocs(favId, 10, 0)
	if err != nil {
		t.Error(err.Error())
	}
	if len(docs) != 3 {
		t.Error("number of docs incorrect")
	}
}
