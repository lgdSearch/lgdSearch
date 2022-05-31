package handler_test

import (
	"errors"
	"lgdSearch/handler"
	"lgdSearch/pkg/db"
	"lgdSearch/pkg/models"
	"testing"

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
	m.Run()
	db.Engine.Where("1 = 1").Delete(&models.User{})
	db.Engine.Where("1 = 1").Delete(&models.Favorite{})
}

func newUserId(name string) uint {
	user := models.User{
		Username: name + "_test",
		Password: "test",
		Favorites: []models.Favorite{
			{DocId: 1},
			{DocId: 2},
			{DocId: 5},
		},
	}
	db.Engine.Create(&user)
	return user.ID
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
	if result.Error != nil || rUser.Nickname != "test" {
		t.Error(result.Error.Error())
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
	err := handler.AppendFavorite(id, 123)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestDeleteFavorite(t *testing.T) {
	id := newUserId("DeleteFavorite")
	err := handler.DeleteFavorite(id, 1)
	if err != nil {
		t.Error(err.Error())
	}
	num := db.Engine.Model(&models.User{Model: gorm.Model{ID: id}}).Association("Favorites").Count()
	if num != 2 {
		t.Error("incorrect number of entries")
	}
}

func TestQueryFavorites(t *testing.T) {
	id := newUserId("QueryFavorites")
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
	id := newUserId("QueryFavorite")
	_, err := handler.QueryFavorite(id, 1)
	if err != nil {
		t.Error(err.Error())
	}
	_, err = handler.QueryFavorite(id, 2)
	if err != nil {
		t.Error(err.Error())
	}
	_, err = handler.QueryFavorite(id, 5)
	if err != nil {
		t.Error(err.Error())
	}
	_, err = handler.QueryFavorite(id, 6)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Error(err.Error())
	}
}
