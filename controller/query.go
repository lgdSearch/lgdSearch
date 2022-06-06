package controller

import (
	"github.com/gin-gonic/gin"
	"lgdSearch/handler"
	"lgdSearch/payloads"
	"lgdSearch/pkg/models"
	"lgdSearch/pkg/weberror"
	"net/http"
)

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