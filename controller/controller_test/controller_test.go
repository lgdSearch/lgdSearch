package controller_test

import (
	"encoding/json"
	"io/ioutil"
	"lgdSearch/payloads"
	"lgdSearch/pkg/db"
	"lgdSearch/pkg/httprequest"
	"lgdSearch/pkg/logger"
	"lgdSearch/pkg/models"
	"lgdSearch/router"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

var engine *gin.Engine

func TestMain(m *testing.M) {
	//初始化日志
	if err := logger.InitLog(logrus.DebugLevel); err != nil {
		panic("log initialization failed")
	}
	//加载配置文件
	//环境变量>配置文件>默认值
	viper.SetConfigFile("test.env")
	viper.AutomaticEnv() //自动匹配环境
	if err := viper.ReadInConfig(); err != nil {
		logger.Logger.Errorf("Fatal error config file:%s", err)
		return
	}
	engine = router.Init()
	db.Init()

	// trie.InitHotSearch("./pkg/data/HotSearch.txt")
	// defer trie.GetHotSearch().Flush("./pkg/data/HotSearch.txt") // flush
	// log.Println("HotSearch Init Success!")

	// trie.InitTrie("./pkg/data/trieData.txt") // 载入 trie
	// defer trie.Tree.FlushIndex("./pkg/data/trieData.txt")
	// log.Println("Trie Init Success!")

	engine.Run(":9090")
	os.Exit(m.Run())
}

func TestRegister(t *testing.T) {
	defer db.Engine.Where("1 = 1").Delete(&models.User{})
	uri := "/register"
	params := map[string]interface{}{
		"username": "testRegister",
		"password": "test",
	}
	w := httprequest.Put("", uri, params, engine)
	result := w.Result()
	defer result.Body.Close()
	body, _ := ioutil.ReadAll(result.Body)
	if w.Code != 204 {
		t.Errorf("code:%d err:%v", w.Code, string(body))
	}
}

func TestLogin(t *testing.T) {
	defer db.Engine.Where("1 = 1").Delete(&models.User{})
	hash, _ := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
	user := models.User{
		Username: "test",
		Password: string(hash),
		Nickname: "游客",
	}
	db.Engine.Create(&user)
	uri := "/login"
	params := map[string]interface{}{
		"username": "test",
		"password": "test",
	}
	w := httprequest.Post("", uri, params, engine)
	result := w.Result()
	defer result.Body.Close()
	body, _ := ioutil.ReadAll(result.Body)
	if w.Code != 200 {
		t.Errorf("code:%d err:%v", w.Code, string(body))
	}
}

func TestLogout(t *testing.T) {
	defer db.Engine.Where("1 = 1").Delete(&models.User{})
	hash, _ := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
	user := models.User{
		Username: "test",
		Password: string(hash),
		Nickname: "游客",
	}
	db.Engine.Create(&user)
	uri := "/login"
	params := map[string]interface{}{
		"username": "test",
		"password": "test",
	}
	w := httprequest.Post("", uri, params, engine)
	r1 := w.Result()
	defer r1.Body.Close()
	body, _ := ioutil.ReadAll(r1.Body)
	if w.Code != 200 {
		t.Errorf("code:%d err:%v", w.Code, string(body))
	}
	var resp payloads.LoginResp
	err := json.Unmarshal(body, &resp)
	if err != nil {
		t.Error(err.Error())
	}
	uri = "/users/logout"
	w = httprequest.Put(resp.Token, uri, nil, engine)
	r2 := w.Result()
	defer r2.Body.Close()
	body, _ = ioutil.ReadAll(r2.Body)
	if w.Code != 200 {
		t.Errorf("code:%d err:%v", w.Code, string(body))
	}
}
