package controller

import (
	"lgdSearch/handler"
	"lgdSearch/payloads"
	"lgdSearch/pkg/models"
	"lgdSearch/pkg/weberror"
	"lgdSearch/pkg/logger"
	"net/http"
	"github.com/gin-gonic/gin"
	jwt "github.com/appleboy/gin-jwt/v2"
)

func Register(c *gin.Context) {
	var req payloads.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Logger.Errorf("[Register] failed to parse request, err: %s", err.Error())
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	_, err := handler.CreateUser(req.Username, req.Password)
	if err != nil {
		logger.Logger.Errorf("[Register] failed to create user, err: %s", err.Error())
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func UpdateProfile(c *gin.Context) {
	var req payloads.UpdateProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Logger.Errorf("[UpdateProfile] failed to parse request, err: %s", err.Error())
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	claims := jwt.ExtractClaims(c)
	err := handler.UpdateUser(claims["user"].(*models.User).ID, req.Nickname)
	if err != nil {
		logger.Logger.Errorf("[UpdateProfile] failed to update user, err: %s", err.Error())
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}