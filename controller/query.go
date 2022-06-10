package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"lgdSearch/handler"
	"lgdSearch/middleware"
	"lgdSearch/payloads"
	"lgdSearch/pkg/models"
	"lgdSearch/pkg/weberror"
	"net/http"
)

// 搜索
// @Tags  search
// @Description
// @Accept   json
// @Produce  json
// @Param     SearchRequest    body      models.SearchRequest  true "searchRequest"
// @Success  200  {object}  payloads.Result{data=models.SearchResult}
// @Failure  400  {object}  weberror.Info  "Bad Request"
// @Failure  404  {object}  weberror.Info  "Not Found"
// @Failure  500  {object}  weberror.Info  "InternalServerError"
// @Router   /query [post]
func Query(c *gin.Context) {
	token, err := middleware.MW.ParseToken(c)
	var userid uint
	if err == nil {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			iUser, _ := claims["user"].(map[string]interface{})
			userid = uint(iUser["ID"].(float64))
		}
	}
	var request = &models.SearchRequest{}
	err = c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	request.Likes = make(map[uint]models.Docs, 0)
	if userid != 0 {
		favorites, _ := handler.QueryFavoritesAndDocs(userid)
		for _, fav := range favorites {
			for _, doc := range fav.Docs {
				request.Likes[doc.DocIndex] = models.Docs{fav.ID, doc.ID}
			}
		}
	}
	r := handler.MultiSearch(request)
	c.JSON(http.StatusOK, payloads.Success(r))
}
