package httprequest

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func Get(token, uri string, bodyParam map[string]interface{}, router *gin.Engine) *httptest.ResponseRecorder {
	// 将参数转化为json比特流
	jsonByte, _ := json.Marshal(bodyParam)
	// 构造get请求，json数据以请求body的形式传递
	req := httptest.NewRequest("GET", uri, bytes.NewReader(jsonByte))
	req.Header.Set("Authorization", "Bearer "+token)
	// 初始化响应
	w := httptest.NewRecorder()
	// 调用相应的handler接口
	router.ServeHTTP(w, req)
	return w
}

func Delete(token, uri string, bodyParam map[string]interface{}, router *gin.Engine) *httptest.ResponseRecorder {
	// 将参数转化为json比特流
	jsonByte, _ := json.Marshal(bodyParam)
	// 构造get请求，json数据以请求body的形式传递
	req := httptest.NewRequest("DELETE", uri, bytes.NewReader(jsonByte))
	req.Header.Set("Authorization", "Bearer "+token)
	// 初始化响应
	w := httptest.NewRecorder()
	// 调用相应的handler接口
	router.ServeHTTP(w, req)
	return w
}

func Post(token, uri string, bodyParam map[string]interface{}, router *gin.Engine) *httptest.ResponseRecorder {
	// 将参数转化为json比特流
	jsonByte, _ := json.Marshal(bodyParam)
	// 构造post请求，json数据以请求body的形式传递
	req := httptest.NewRequest("POST", uri, bytes.NewReader(jsonByte))
	req.Header.Set("Authorization", "Bearer "+token)
	// 初始化响应
	w := httptest.NewRecorder()
	// 调用相应的handler接口
	router.ServeHTTP(w, req)
	return w
}

func Patch(token, uri string, bodyParam map[string]interface{}, router *gin.Engine) *httptest.ResponseRecorder {
	// 将参数转化为json比特流
	jsonByte, _ := json.Marshal(bodyParam)
	// 构造PATCH请求，json数据以请求body的形式传递
	req := httptest.NewRequest("PATCH", uri, bytes.NewReader(jsonByte))
	req.Header.Set("Authorization", "Bearer "+token)
	// 初始化响应
	w := httptest.NewRecorder()
	// 调用相应的handler接口
	router.ServeHTTP(w, req)
	return w
}

func Put(token, uri string, bodyParam interface{}, router *gin.Engine) *httptest.ResponseRecorder {
	// 将参数转化为json比特流
	jsonByte, _ := json.Marshal(bodyParam)
	// 构造put请求，json数据以请求body的形式传递
	req := httptest.NewRequest("PUT", uri, bytes.NewReader(jsonByte))
	req.Header.Set("Authorization", "Bearer "+token)
	// 初始化响应
	w := httptest.NewRecorder()

	// 调用相应的handler接口
	router.ServeHTTP(w, req)
	return w
}

