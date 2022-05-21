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

func AddFavorite(c *gin.Context) {
	var req payloads.AddFavoriteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Logger.Errorf("[AddFavorite] failed to parse request, err: %s", err.Error())
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	claims := jwt.ExtractClaims(c)
	err := handler.AppendFavorite(claims["user"].(*models.User).ID, req.DocId)
	if err != nil {
		logger.Logger.Errorf("[AddFavorite] failed to addend favorite, err: %s", err.Error())
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func DeleteFavorite(c *gin.Context) {
	var req payloads.DeleteFavoriteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Logger.Errorf("[DeleteFavorite] failed to parse request, err: %s", err.Error())
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	claims := jwt.ExtractClaims(c)
	err := handler.DeleteFavorite(claims["user"].(*models.User).ID, req.DocId)
	if err != nil {
		logger.Logger.Errorf("[DeleteFavorite] failed to delete favorite, err: %s", err.Error())
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func GetFavorites(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	favorites, err := handler.QueryFavorites(claims["user"].(*models.User).ID)
	if err != nil {
		logger.Logger.Errorf("[DeleteFavorite] failed to query favorites, err: %s", err.Error())
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	resps := make([]payloads.GetFavoritesResp, 0, len(favorites))
	for _, v := range favorites {
		resps = append(resps, payloads.GetFavoritesResp{
			DocId: v.DocId,
			Summary: v.Summary,
		})
	}
	c.JSON(http.StatusOK, resps)
}