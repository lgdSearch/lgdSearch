package controller

import (
	"github.com/gin-gonic/gin"
	"lgdSearch/payloads"
	"lgdSearch/pkg/logger"
	"lgdSearch/pkg/weberror"
	"lgdSearch/handler"
	"net/http"
)

func SayHello(c *gin.Context) {
	var req payloads.SayHelloReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Logger.Errorf("[SayHello] failed to parse request, err: %s", err.Error())
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	text, err := handler.QueryText(req.Text)
	if err != nil {
		logger.Logger.Errorf("[SayHello] failed to get text, err: %s", err.Error())
		c.JSON(http.StatusInternalServerError, weberror.Info{Error: http.StatusText(http.StatusInternalServerError)})
		return
	}
	c.JSON(http.StatusOK, payloads.SayHelloResp{
		Text: text,
	})
}