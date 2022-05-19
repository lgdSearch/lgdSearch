package controller

import (
	"github.com/gin-gonic/gin"
	"lgdSearch/handler"
	"lgdSearch/pkg/errors"
	"lgdSearch/pkg/logger"
	"net/http"
)

func GetRelatedSearch(c *gin.Context) {
	query := c.Param("text")
	text, err := handler.RelatedQuery(query)
	if err != nil {
		logger.Logger.Errorf("[SayHello] failed to get text, err: %s", err.Error())
		c.JSON(http.StatusInternalServerError, errors.Error{Error: http.StatusText(http.StatusInternalServerError)})
		return
	}
	c.JSON(http.StatusOK, text)
}
