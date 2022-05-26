package controller

import (
	"github.com/gin-gonic/gin"
	"lgdSearch/handler"
	"lgdSearch/pkg/logger"
	"lgdSearch/pkg/weberror"
	"net/http"
)

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
