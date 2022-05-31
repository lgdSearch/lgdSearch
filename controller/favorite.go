package controller

import (
	"errors"
	"lgdSearch/handler"
	"lgdSearch/payloads"
	"lgdSearch/pkg/logger"
	"lgdSearch/pkg/weberror"
	"lgdSearch/pkg/extractclaims"
	"net/http"
	"strconv"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 添加收藏
// @Tags favorite
// @Description
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string          true       "userToken"
// @Param        doc_id         path      uint            true       "doc_id in leveldb"
// @Success      204
// @Failure      400            {object}  weberror.Info              "Bad Request"
// @Failure      404            {object}  weberror.Info              "Not Found"
// @Failure      500            {object}  weberror.Info              "InternalServerError"
// @Router       /users/favorites/{doc_id} [put]
// @Security     Token
func AddFavorite(c *gin.Context) {
	docId, err := strconv.ParseUint(c.Param("doc_id"), 10, 32)
	if err != nil {
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
	_, err = handler.QueryFavorite(user.ID, uint(docId))
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		if err != nil {
			logger.Logger.Errorf("[AddFavorite] failed to parse request, err: %s", err.Error())
			c.JSON(http.StatusInternalServerError,  weberror.Info{Error: http.StatusText(http.StatusInternalServerError)})
			return
		}
		logger.Logger.Errorf("[AddFavorite] failed to parse request, err: %s", "duplicate document")
		c.JSON(http.StatusInternalServerError, weberror.Info{Error: "duplicate document"})
		return
	}
	err = handler.AppendFavorite(user.ID, uint(docId))
	if err != nil {
		logger.Logger.Errorf("[AddFavorite] failed to addend favorite, err: %s", err.Error())
		c.JSON(http.StatusInternalServerError,  weberror.Info{Error: http.StatusText(http.StatusInternalServerError)})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// 取消收藏
// @Tags favorite
// @Description
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string          true       "userToken"
// @Param        doc_id         path      uint            true       "doc_id in leveldb"
// @Success      204
// @Failure      400            {object}  weberror.Info              "Bad Request"
// @Failure      404            {object}  weberror.Info              "Not Found"
// @Failure      500            {object}  weberror.Info              "InternalServerError"
// @Router       /users/favorites/{doc_id} [delete]
// @Security     Token
func DeleteFavorite(c *gin.Context) {
	docId, err := strconv.ParseUint(c.Param("doc_id"), 10, 32)
	if err != nil {
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
	err = handler.DeleteFavorite(user.ID, uint(docId))
	if err != nil {
		logger.Logger.Errorf("[DeleteFavorite] failed to delete favorite, err: %s", err.Error())
		c.JSON(http.StatusInternalServerError,  weberror.Info{Error: http.StatusText(http.StatusInternalServerError)})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// 获取全部收藏
// @Tags favorite
// @Description
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string          true       "userToken"
// @Param        limit          query     uint            false      "greater than 0"
// @Param        offset         query     uint            false		 "greater than -1"
// @Success      200            {object}  payloads.GetFavoritesResp
// @Failure      400            {object}  weberror.Info              "Bad Request"
// @Failure      404            {object}  weberror.Info              "Not Found"
// @Failure      500            {object}  weberror.Info              "InternalServerError"
// @Router       /users/favorites [get]
// @Security     Token
func GetFavorites(c *gin.Context) {
	limit, err := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 32)
	if err != nil {
		logger.Logger.Errorf("[DeleteFavorite] failed to query favorites, err: %s", err.Error())
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	offset, err := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 32)
	if err != nil {
		logger.Logger.Errorf("[DeleteFavorite] failed to query favorites, err: %s", err.Error())
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	if limit < 1 || offset < 0 {
		logger.Logger.Errorf("[DeleteFavorite] failed to query favorites,err: %s", "Parameter is out of range")
		c.JSON(http.StatusBadRequest, weberror.Info{Error: "Parameter is out of range"})
		return
	}
	user := extractclaims.ToUser(jwt.ExtractClaims(c))
	if user == nil {
		logger.Logger.Errorf("[DeleteFavorite] failed to query favorites, err: %s", "failed to extract user info")
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	favorites, err := handler.QueryFavorites(user.ID, uint(limit), uint(offset))
	if err != nil {
		logger.Logger.Errorf("[DeleteFavorite] failed to query favorites, err: %s", err.Error())
		c.JSON(http.StatusInternalServerError,  weberror.Info{Error: http.StatusText(http.StatusInternalServerError)})
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