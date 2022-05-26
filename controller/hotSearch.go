package controller

import (
	"github.com/gin-gonic/gin"
	"lgdSearch/handler"
	"lgdSearch/payloads"
	"net/http"
)

func HotSearch(c *gin.Context) {
	r := handler.HotSearch()
	c.JSON(http.StatusOK, payloads.Success(r))
}
