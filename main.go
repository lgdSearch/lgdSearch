package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"lgdSearch/pkg/logger"
	"lgdSearch/router"
	"lgdSearch/pkg/db"
)

func main() {
	//初始化日志
	if err := logger.InitLog(logrus.DebugLevel); err != nil {
		panic("log initialization failed")
	}
	//加载配置文件
	//环境变量>配置文件>默认值
	viper.SetConfigFile(".env")
	viper.AutomaticEnv() //自动匹配环境
	if err := viper.ReadInConfig(); err != nil {
		logger.Logger.Errorf("Fatal error config file:%s", err)
		return
	}
	engine := router.Init()
	db.Init()
	go engine.Run(":9090")
}
