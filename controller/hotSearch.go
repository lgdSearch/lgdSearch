package controller

import (
	"github.com/gin-gonic/gin"
	"lgdSearch/handler"
	"lgdSearch/payloads"
	"net/http"
)

// 热搜
// @Tags  search
// @Description
// @Accept   json
// @Produce  json
// @Success  200  {object}  payloads.Result{data=[]trie.HotSearchMessage}
// @Failure  400  {object}  weberror.Info  "Bad Request"
// @Failure  404  {object}  weberror.Info  "Not Found"
// @Failure  500  {object}  weberror.Info  "InternalServerError"
// @Router   /query/hotSearch [get]
func HotSearch(c *gin.Context) {
	r := handler.HotSearch()
	c.JSON(http.StatusOK, payloads.Success(r))
}
