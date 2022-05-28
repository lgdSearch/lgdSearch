package controller

import (
	"errors"
	"lgdSearch/handler"
	"lgdSearch/payloads"
	"lgdSearch/pkg/logger"
	"lgdSearch/pkg/models"
	"lgdSearch/pkg/weberror"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(c *gin.Context) {
	var req payloads.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Logger.Errorf("[Register] failed to parse request, err: %s", err.Error())
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	_, err := handler.QueryUser(0, req.Username)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		if err != nil {
			logger.Logger.Errorf("[Register] failed to create user, err: %s", err.Error())
			c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
			return
		}
		logger.Logger.Errorf("[Register] failed to create user, err: %s", "duplicate username")
		c.JSON(http.StatusBadRequest, weberror.Info{Error: "duplicate username"})
		return
	}
	_, err = handler.CreateUser(req.Username, req.Password)
	if err != nil {
		logger.Logger.Errorf("[Register] failed to create user, err: %s", err.Error())
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func UpdateNickname(c *gin.Context) {
	var req payloads.UpdateProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Logger.Errorf("[UpdateProfile] failed to parse request, err: %s", err.Error())
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	claims := jwt.ExtractClaims(c)
	err := handler.UpdateUserNickname(claims["user"].(*models.User).ID, req.Nickname)
	if err != nil {
		logger.Logger.Errorf("[UpdateProfile] failed to update user, err: %s", err.Error())
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func DeleteAccount(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	err := handler.DeleteUser(claims["user"].(*models.User).ID)
	if err != nil {
		logger.Logger.Errorf("[DeleteAccount] failed to delete user, err: %s", err.Error())
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func GetProfile(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	user, err := handler.QueryUser(claims["user"].(*models.User).ID, "")
	if err != nil {
		logger.Logger.Errorf("[GetProfile] failed to query user, err: %s", err.Error())
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	c.JSON(http.StatusOK, &payloads.GetProfileResp{
		Username: user.Username,
		Nickname: user.Nickname,
	})
}