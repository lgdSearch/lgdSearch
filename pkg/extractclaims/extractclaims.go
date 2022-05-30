package extractclaims

import (
	"lgdSearch/pkg/models"
	jwt "github.com/appleboy/gin-jwt/v2"
)

func ToUser (claims jwt.MapClaims) *models.User{
	iUser, ok := claims["user"].(map[string]interface{})
	if !ok {
		return nil
	}
	
	id, ok := iUser["ID"].(float64)
	if !ok {
		return nil
	}
	username, ok := iUser["Username"].(string)
	if !ok {
		return nil
	}
	nickname, ok := iUser["Nickname"].(string)
	if !ok {
		return nil
	}
	user := models.User{}
	user.ID = uint(id)
	user.Username = username
	user.Nickname = nickname
	return &user
}