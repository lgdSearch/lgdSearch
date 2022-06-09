package middleware

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"lgdSearch/handler"
	"lgdSearch/payloads"
	"lgdSearch/pkg/models"
	"lgdSearch/pkg/weberror"
	"net/http"
)

const identityKey = "user"

var MW *jwt.GinJWTMiddleware = nil

func GetJWTMiddle() *jwt.GinJWTMiddleware {
	if MW == nil {
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
				if err := c.ShouldBindJSON(&req); err != nil {
					return "", jwt.ErrMissingLoginValues
				}
				username := req.Username
				password := req.Password
				user, err := handler.QueryUser(0, username)
				if (err == nil) && (bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil) {
					return user, nil
				}
				return nil, jwt.ErrFailedAuthentication
			},
			Unauthorized: func(c *gin.Context, code int, message string) {
				c.JSON(code, weberror.Info{Error: message})
			},
			LogoutResponse: func(c *gin.Context, code int) {
				c.JSON(http.StatusNoContent, nil)
			},
		})
		if err != nil {
			panic("JWTMiddle initialization failed")
		}
		MW = authMiddleware
	}
	return MW
}
