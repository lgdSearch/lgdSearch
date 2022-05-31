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
	m.Run()
	db.Engine.Where("1 = 1").Delete(&models.User{})
	db.Engine.Where("1 = 1").Delete(&models.Favorite{})
}

func newUserToken(name string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
	user := models.User{
		Username: name + "_test",
		Password: string(hash),
		Nickname: "游客",
		Favorites: []models.Favorite{
			{DocId: 1},
			{DocId: 2},
			{DocId: 5},
		},
	}
	db.Engine.Create(&user)
	uri := "/login"
	params := map[string]interface{}{
		"username": name + "_test",
		"password": "test",
	}
	w := httprequest.Post("", uri, params, engine)
	r1 := w.Result()
	defer r1.Body.Close()
	if w.Code != 200 {
		return ""
	}
	body, _ := ioutil.ReadAll(r1.Body)
	var resp payloads.LoginResp
	json.Unmarshal(body, &resp)
	return resp.Token
}

func TestRegister(t *testing.T) {
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
	token := newUserToken("Logout")
	uri := "/users/logout"
	w := httprequest.Put(token, uri, nil, engine)
	r := w.Result()
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	if w.Code != 204 {
		t.Errorf("code:%d err:%v", w.Code, string(body))
	}
}

func TestUpdateNickname(t *testing.T) {
	token := newUserToken("UpdateNickname")
	uri := "/users/nickname"
	params := map[string]interface{}{
		"nickname": "lgdSearch",
	}
	w := httprequest.Patch(token, uri, params, engine)
	r := w.Result()
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	if w.Code != 204 {
		t.Errorf("code:%d err:%v", w.Code, string(body))
	}
}

func TestDeleteAccount(t *testing.T) {
	token := newUserToken("DeleteAccount")
	uri := "/users"
	w := httprequest.Delete(token, uri, nil, engine)
	r := w.Result()
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	if w.Code != 204 {
		t.Errorf("code:%d err:%v", w.Code, string(body))
	}
}

func TestGetProfile(t *testing.T) {
	token := newUserToken("GetProfile")
	uri := "/users/profile"
	w := httprequest.Get(token, uri, nil, engine)
	r := w.Result()
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	if w.Code != 200 {
		t.Errorf("code:%d err:%v", w.Code, string(body))
	}
}

func TestAddFavorite(t *testing.T) {
	token := newUserToken("AddFavorite")
	uri := "/users/favorites/123"
	w := httprequest.Put(token, uri, nil, engine)
	r := w.Result()
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	if w.Code != 204 {
		t.Errorf("code:%d err:%v", w.Code, string(body))
	}
}

func TestDeleteFavorite(t *testing.T) {
	token := newUserToken("DeleteFavorite")
	uri := "/users/favorites/1"
	w := httprequest.Delete(token, uri, nil, engine)
	r := w.Result()
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	if w.Code != 204 {
		t.Errorf("code:%d err:%v", w.Code, string(body))
	}
}

func TestGetFavorites(t *testing.T) {
	token := newUserToken("GetFavorites")
	uri := "/users/favorites"
	w := httprequest.Get(token, uri, nil, engine)
	r := w.Result()
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	if w.Code != 200 {
		t.Errorf("code:%d err:%v", w.Code, string(body))
	}
}
