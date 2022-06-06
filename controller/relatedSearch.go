package controller

import (
	"github.com/gin-gonic/gin"
	"lgdSearch/handler"
	"lgdSearch/pkg/logger"
	"lgdSearch/pkg/weberror"
	"net/http"
)

// 相关搜索
// @Tags  search
// @Description
// @Accept   json
// @Produce  json
// @Param     text         path    string    true  "text"
// @Success  200  {object}  handler.PageData{data=handler.PageInfo}
// @Failure  400  {object}  weberror.Info  "Bad Request"
// @Failure  404  {object}  weberror.Info  "Not Found"
// @Failure  500  {object}  weberror.Info  "InternalServerError"
// @Router   /book/{text} [get]
func GetRelatedSearch(c *gin.Context) {
	query := c.Param("text")
	text, err := handler.RelatedQuery(query)
	if err != nil {
		logger.Logger.Errorf("[SayHello] failed to get text, err: %s", err.Error())
		c.JSON(http.StatusInternalServerError, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	c.JSON(http.StatusOK, text)
}
