package middleware

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"lgdSearch/pkg/models"
	"lgdSearch/payloads"
	"lgdSearch/pkg/weberror"
	"lgdSearch/handler"
	"github.com/gin-gonic/gin"
)

const identityKey = "user"

func GetJWTMiddle() *jwt.GinJWTMiddleware {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "lgdSearch",
		Key:         []byte("bu ban meng ma"),
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.User); ok {
				return jwt.MapClaims{
					identityKey: v,
				}
			}
			return jwt.MapClaims{}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var req payloads.LoginReq
			if err := c.ShouldBind(&req); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			username := req.Username
			password := req.Password
			user, err := handler.QueryUser(0, username)
			if (err == nil) && (username == user.Username && password == user.Password) {
				return user, nil
			}
			return nil, jwt.ErrFailedAuthentication
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, weberror.Info{Error: message})
		},
	})
	if err != nil {
		panic("JWTMiddle initialization failed")
	}
	return authMiddleware
}
