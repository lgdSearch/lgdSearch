package controller

import (
	"errors"
	"lgdSearch/handler"
	"lgdSearch/payloads"
	"lgdSearch/pkg/logger"
	"lgdSearch/pkg/weberror"
	"lgdSearch/pkg/extractclaims"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 用户注册
// @Tags user
// @Description
// @Accept       json
// @Produce      json
// @Param        RegisterReq    body      payloads.RegisterReq  true  "username and pwd"
// @Success      204
// @Failure      400            {object}  weberror.Info               "Bad Request"
// @Failure      404            {object}  weberror.Info               "Not Found"
// @Failure      500            {object}  weberror.Info               "InternalServerError"
// @Router       /register [put]
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
			c.JSON(http.StatusInternalServerError,  weberror.Info{Error: http.StatusText(http.StatusInternalServerError)})
			return
		}
		logger.Logger.Errorf("[Register] failed to create user, err: %s", "duplicate username")
		c.JSON(http.StatusInternalServerError,  weberror.Info{Error: "duplicate username"})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Logger.Errorf("[Register] failed to hash password, err: %s", err.Error())
		c.JSON(http.StatusInternalServerError,  weberror.Info{Error: http.StatusText(http.StatusInternalServerError)})
		return
	}
	req.Password = string(hash)
	_, err = handler.CreateUser(req.Username, req.Password)
	if err != nil {
		logger.Logger.Errorf("[Register] failed to create user, err: %s", err.Error())
		c.JSON(http.StatusInternalServerError,  weberror.Info{Error: http.StatusText(http.StatusInternalServerError)})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// 修改昵称
// @Tags user
// @Description
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                     true "userToken"
// @Param        UpdateNickname body      payloads.UpdateProfileReq  true "nickname"
// @Success      204
// @Failure      400            {object}  weberror.Info                   "Bad Request"
// @Failure      404            {object}  weberror.Info                   "Not Found"
// @Failure      500            {object}  weberror.Info                   "InternalServerError"
// @Router       /users/nickname [patch]
// @Security     Token
func UpdateNickname(c *gin.Context) {
	var req payloads.UpdateProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Logger.Errorf("[UpdateProfile] failed to parse request, err: %s", err.Error())
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	user := extractclaims.ToUser(jwt.ExtractClaims(c))
	if user == nil {
		logger.Logger.Errorf("[UpdateProfile] failed to update user, err: %s", "failed to extract user info")
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	err := handler.UpdateUserNickname(user.ID, req.Nickname)
	if err != nil {
		logger.Logger.Errorf("[UpdateProfile] failed to update user, err: %s", err.Error())
		c.JSON(http.StatusInternalServerError,  weberror.Info{Error: http.StatusText(http.StatusInternalServerError)})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// 删除账户
// @Tags user
// @Description
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string         true  "userToken"
// @Success      204
// @Failure      400            {object}  weberror.Info        "Bad Request"
// @Failure      404            {object}  weberror.Info        "Not Found"
// @Failure      500            {object}  weberror.Info        "InternalServerError"
// @Router       /users [delete]
// @Security     Token
func DeleteAccount(c *gin.Context) {
	user := extractclaims.ToUser(jwt.ExtractClaims(c))
	if user == nil {
		logger.Logger.Errorf("[DeleteAccount] failed to delete user, err: %s", "failed to extract user info")
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	err := handler.DeleteUser(user.ID)
	if err != nil {
		logger.Logger.Errorf("[DeleteAccount] failed to delete user, err: %s", err.Error())
		c.JSON(http.StatusInternalServerError,  weberror.Info{Error: http.StatusText(http.StatusInternalServerError)})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// 获取个人资料
// @Tags user
// @Description
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                  true "userToken"
// @Success      200            {object}  payloads.GetProfileResp
// @Failure      400            {object}  weberror.Info                "Bad Request"
// @Failure      404            {object}  weberror.Info                "Not Found"
// @Failure      500            {object}  weberror.Info                "InternalServerError"
// @Router       /users/profile [get]
// @Security     Token
func GetProfile(c *gin.Context) {
	user := extractclaims.ToUser(jwt.ExtractClaims(c))
	if user == nil {
		logger.Logger.Errorf("[GetProfile] failed to query user, err: %s", "failed to extract user info")
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	user, err := handler.QueryUser(user.ID, "")
	if err != nil {
		logger.Logger.Errorf("[GetProfile] failed to query user, err: %s", err.Error())
		c.JSON(http.StatusInternalServerError,  weberror.Info{Error: http.StatusText(http.StatusInternalServerError)})
		return
	}
	c.JSON(http.StatusOK, &payloads.GetProfileResp{
		Username: user.Username,
		Nickname: user.Nickname,
	})
}

// 登录
// @Tags user
// @Description
// @Accept       json
// @Produce      json
// @Param        LoginReq       body      payloads.LoginReq    true  "username and pwd"
// @Success      200            {object}  payloads.LoginResp
// @Failure      400            {object}  weberror.Info              "Bad Request"
// @Failure      404            {object}  weberror.Info              "Not Found"
// @Failure      500            {object}  weberror.Info              "InternalServerError"
// @Router       /login [post]
func Login(c *gin.Context) {
	//生成文档用，真实接口由gin-jwt生成
}

// 登出
// @Tags user
// @Description
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string         true  "userToken"
// @Success      204
// @Failure      400            {object}  weberror.Info        "Bad Request"
// @Failure      404            {object}  weberror.Info        "Not Found"
// @Failure      500            {object}  weberror.Info        "InternalServerError"
// @Router       /users/logout [delete]
// @Security     Token
func Logout(c *gin.Context) {
	//生成文档用，真实接口由gin-jwt生成
}