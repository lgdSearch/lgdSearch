package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"lgdSearch/pkg"
	"lgdSearch/pkg/db"
	"lgdSearch/pkg/logger"
	"lgdSearch/pkg/trie"
	"lgdSearch/router"
	"time"
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

	PkgEngine := pkg.Engine{IndexPath: "./pkg/data"}
	PkgEngine.Init()
	defer PkgEngine.Close()

	trie.InitHotSearch("./pkg/data/HotSearch.txt")
	defer trie.GetHotSearch().Flush("./pkg/data/HotSearch.txt") // flush

	trie.InitTrie("./pkg/data/trieData.txt") // 载入 trie
	defer trie.Tree.FlushIndex("./pkg/data/trieData.txt")

	pkg.Set(&PkgEngine)

	go engine.Run(":9090")

	time.Sleep(time.Second * 60) // 运行时间
}
