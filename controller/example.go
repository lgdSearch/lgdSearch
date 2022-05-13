package controller

import (
	"github.com/gin-gonic/gin"
	"lgdSearch/payloads"
	"lgdSearch/pkg/logger"
	"lgdSearch/pkg/errors"
	"lgdSearch/handler"
	"net/http"
)

func SayHello(c *gin.Context) {
	var req payloads.SayHelloReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Logger.Errorf("[SayHello] failed to parse request, err: %s", err.Error())
		c.JSON(http.StatusBadRequest, errors.Error{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	text, err := handler.RetrieveText(req.Text)
	if err != nil {
		logger.Logger.Errorf("[SayHello] failed to get text, err: %s", err.Error())
		c.JSON(http.StatusInternalServerError, errors.Error{Error: http.StatusText(http.StatusInternalServerError)})
		return
	}
	c.JSON(http.StatusOK, payloads.SayHelloResp{
		Text: text,
	})
}