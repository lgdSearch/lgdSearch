package db

import (
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
	"github.com/spf13/viper"
	"lgdSearch/pkg/models"
)

var Engine *gorm.DB

func Init() {
	var err error
	Engine, err = gorm.Open(mysql.Open(viper.GetString("mysql_dsn")))
	if err != nil {
		panic("failed to connect database")
	}
	Engine.AutoMigrate(&models.User{})
	Engine.AutoMigrate(&models.Favorite{})
	Engine.AutoMigrate(&models.Doc{})
}

