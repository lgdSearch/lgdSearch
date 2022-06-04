package controller

import (
	"github.com/gin-gonic/gin"
	"lgdSearch/handler"
	"lgdSearch/payloads"
	"lgdSearch/pkg/models"
	"lgdSearch/pkg/weberror"
	"net/http"
)

func QueryPicture(c *gin.Context) {
	var request = &models.SearchRequest{}
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	r := handler.MultiSearchPicture(request)
	c.JSON(http.StatusOK, payloads.Success(r))
}
