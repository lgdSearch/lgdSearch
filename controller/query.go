package controller

import (
	"github.com/gin-gonic/gin"
	"lgdSearch/handler"
	"lgdSearch/payloads"
	"lgdSearch/pkg/models"
	"lgdSearch/pkg/weberror"
	"net/http"
)

// 搜索
// @Tags  search
// @Description
// @Accept   json
// @Produce  json
// @Param     SearchRequest    body      models.SearchRequest  true "searchRequest"
// @Success  200  {object}  payloads.Result{data=models.SearchResult}
// @Failure  400  {object}  weberror.Info  "Bad Request"
// @Failure  404  {object}  weberror.Info  "Not Found"
// @Failure  500  {object}  weberror.Info  "InternalServerError"
// @Router   /query [post]
func Query(c *gin.Context) {
	var request = &models.SearchRequest{}
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	r := handler.MultiSearch(request)
	c.JSON(http.StatusOK, payloads.Success(r))
}
