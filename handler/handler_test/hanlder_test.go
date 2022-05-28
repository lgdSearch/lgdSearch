package handler_test

import (
	"testing"
	"lgdSearch/pkg/db"
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
	"github.com/spf13/viper"
	"lgdSearch/pkg/models"
	"lgdSearch/handler"
	"os"
	"errors"
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
	os.Exit(m.Run())
}

func TestCreateUser(t *testing.T) {
	defer db.Engine.Where("1 = 1").Delete(&models.User{})
	_, err := handler.CreateUser("test", "123456")
	if err != nil {
		t.Error(err.Error())
	}
}

func TestQueryUser(t *testing.T) {
	defer db.Engine.Where("1 = 1").Delete(&models.User{})
	user := models.User{
		Username: "test",
		Password: "123456",
	}
	db.Engine.Create(&user)
	_, err := handler.QueryUser(user.ID, "")
	if err != nil {
		t.Error(err.Error())
	}
	_, err = handler.QueryUser(0, "test")
	if err != nil {
		t.Error(err.Error())
	}
}

func TestUpdateUserNickname(t *testing.T) {
	defer db.Engine.Where("1 = 1").Delete(&models.User{})
	user := models.User{
		Username: "test",
		Password: "123456",
	}
	db.Engine.Create(&user)
	defer db.Engine.Where("1 = 1").Delete(&models.User{})
	err := handler.UpdateUserNickname(user.ID, "test")
	if err != nil {
		t.Error(err.Error())
	}
	rUser := models.User{}
	result := db.Engine.First(&rUser, user.ID)
	if result.Error != nil || rUser.Nickname != "test" {
		t.Error(result.Error.Error())
	}
}

func TestDeleteUser(t *testing.T) {
	defer db.Engine.Where("1 = 1").Delete(&models.User{})
	user := models.User{
		Username: "test",
		Password: "123456",
	}
	db.Engine.Create(&user)
	err := handler.DeleteUser(user.ID)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestAppendFavorite (t *testing.T) {
	defer db.Engine.Where("1 = 1").Delete(&models.User{})
	user := models.User{
		Username: "test",
		Password: "123456",
	}
	db.Engine.Create(&user)
	err := handler.AppendFavorite(user.ID, 1)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestDeleteFavorite (t *testing.T) {
	defer db.Engine.Where("1 = 1").Delete(&models.User{})
	user := models.User{
		Username: "test",
		Password: "123456",
		Favorites: []models.Favorite{
			{DocId: 1},
		},
	}
	db.Engine.Create(&user)
	err := handler.DeleteFavorite(user.ID, 1)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestQueryFavorites (t *testing.T) {
	defer db.Engine.Where("1 = 1").Delete(&models.User{})
	user := models.User{
		Username: "test",
		Password: "123456",
		Favorites: []models.Favorite{
			{DocId: 1},
			{DocId: 2},
			{DocId: 5},
		},
	}
	db.Engine.Create(&user)
	favorites, err := handler.QueryFavorites(user.ID)
	if err != nil || len(favorites) != 3{
		t.Error(err.Error())
	}
}

func TestQueryFavorite (t *testing.T) {
	defer db.Engine.Where("1 = 1").Delete(&models.User{})
	user := models.User{
		Username: "test",
		Password: "123456",
		Favorites: []models.Favorite{
			{DocId: 1},
			{DocId: 2},
			{DocId: 5},
		},
	}
	db.Engine.Create(&user)
	_, err := handler.QueryFavorite(user.ID, 1)
	if err != nil{
		t.Error(err.Error())
	}
	_, err = handler.QueryFavorite(user.ID, 2)
	if err != nil{
		t.Error(err.Error())
	}
	_, err = handler.QueryFavorite(user.ID, 5)
	if err != nil{
		t.Error(err.Error())
	}
	_, err = handler.QueryFavorite(user.ID, 6)
	if !errors.Is(err, gorm.ErrRecordNotFound){
		t.Error(err.Error())
	}
}