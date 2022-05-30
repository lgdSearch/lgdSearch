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
	"gorm.io/gorm"
)

func AddFavorite(c *gin.Context) {
	var req payloads.AddFavoriteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Logger.Errorf("[AddFavorite] failed to parse request, err: %s", err.Error())
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	user := extractclaims.ToUser(jwt.ExtractClaims(c))
	if user == nil {
		logger.Logger.Errorf("[AddFavorite] failed to parse request, err: %s", "failed to extract user info")
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	_, err := handler.QueryFavorite(user.ID, req.DocId)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		if err != nil {
			logger.Logger.Errorf("[AddFavorite] failed to parse request, err: %s", err.Error())
			c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
			return
		}
		logger.Logger.Errorf("[AddFavorite] failed to parse request, err: %s", "duplicate document")
		c.JSON(http.StatusBadRequest, weberror.Info{Error: "duplicate document"})
		return
	}
	err = handler.AppendFavorite(user.ID, req.DocId)
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
	user := extractclaims.ToUser(jwt.ExtractClaims(c))
	if user == nil {
		logger.Logger.Errorf("[DeleteFavorite] failed to parse request, err: %s", "failed to extract user info")
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	err := handler.DeleteFavorite(user.ID, req.DocId)
	if err != nil {
		logger.Logger.Errorf("[DeleteFavorite] failed to delete favorite, err: %s", err.Error())
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func GetFavorites(c *gin.Context) {
	user := extractclaims.ToUser(jwt.ExtractClaims(c))
	if user == nil {
		logger.Logger.Errorf("[DeleteFavorite] failed to query favorites, err: %s", "failed to extract user info")
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	favorites, err := handler.QueryFavorites(user.ID)
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