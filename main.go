package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"lgdSearch/pkg"
	"lgdSearch/pkg/db"
	"lgdSearch/pkg/logger"
	"lgdSearch/pkg/trie"
	"lgdSearch/router"
	"log"
	"os"
	"os/signal"
)

// @title           lgdSearch API
// @version         1.0
// @description     This is a simple search engine.

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:9090

// @securityDefinitions.apikey  Token
// @in                          header
// @name                        Authorization
// @description					should be set with extra string "Bearer " before it, sample: "Authorization:Bearer XXXXXXXXXXX(token)"
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

	PkgEngine := pkg.Engine{IndexPath: pkg.DefaultIndexPath()}
	PkgEngine.Init()
	defer PkgEngine.Close()

	trie.InitHotSearch(pkg.DefaultHotSearchPath())
	defer trie.GetHotSearch().Flush(pkg.DefaultHotSearchPath())

	trie.InitTrie(pkg.DefaultTriePath())
	defer trie.Tree.FlushIndex(pkg.DefaultTriePath())

	pkg.Set(&PkgEngine)

	go engine.Run(":9090")

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println("程序即将结束")
	trie.Tree.FlushIndex(pkg.DefaultTriePath())
	trie.GetHotSearch().Flush(pkg.DefaultHotSearchPath())
	PkgEngine.Close()

	log.Fatal("program interrupted")
}
