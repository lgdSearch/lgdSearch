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

// 添加收藏夹
// @Tags favorite
// @Description
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                   true  "userToken"
// @Param        AddFavoriteReq body      payloads.AddFavoriteReq  true  "name"
// @Success      201
// @Failure      400            {object}  weberror.Info                  "Bad Request"
// @Failure      404            {object}  weberror.Info                  "Not Found"
// @Failure      500            {object}  weberror.Info                  "InternalServerError"
// @Router       /users/favorites [put]
// @Security     Token
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

	_, err := handler.QueryFavorite(0, req.Name)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		if err != nil {
			logger.Logger.Errorf("[AddFavorite] failed to query favorite, err: %s", err.Error())
			c.JSON(http.StatusInternalServerError,  weberror.Info{Error: http.StatusText(http.StatusInternalServerError)})
			return
		}
		logger.Logger.Errorf("[AddFavorite] err: %s", "duplicate name")
		c.JSON(http.StatusInternalServerError, weberror.Info{Error: "duplicate name"})
		return
	}
	_, err = handler.AppendFavorite(user.ID, req.Name)
	if err != nil {
		logger.Logger.Errorf("[AddFavorite] failed to add favorite, err: %s", err.Error())
		c.JSON(http.StatusInternalServerError,  weberror.Info{Error: http.StatusText(http.StatusInternalServerError)})
		return
	}
	c.JSON(http.StatusCreated, nil)
}

// 更新收藏夹名字
// @Tags favorite
// @Description
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                          true  "userToken"
// @Param        fav_id         path      uint                            true  "fav_id"
// @Param        AddFavoriteReq body      payloads.UpdateFavoriteNameReq  true  "name"
// @Success      204
// @Failure      400            {object}  weberror.Info                         "Bad Request"
// @Failure      404            {object}  weberror.Info                         "Not Found"
// @Failure      500            {object}  weberror.Info                         "InternalServerError"
// @Router       /users/favorites/{fav_id} [patch]
// @Security     Token
func UpdateFavoriteName(c *gin.Context) {
	var req payloads.UpdateFavoriteNameReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Logger.Errorf("[UpdateFavoriteName] failed to parse request, err: %s", err.Error())
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	favId, err := strconv.ParseUint(c.Param("fav_id"), 10, 32)
	if err != nil {
		logger.Logger.Errorf("[UpdateFavoriteName] failed to parse request, err: %s", err.Error())
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	user := extractclaims.ToUser(jwt.ExtractClaims(c))
	if user == nil {
		logger.Logger.Errorf("[UpdateFavoriteName] failed to parse request, err: %s", "failed to extract user info")
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	err = handler.UpdateFavoriteName(user.ID, uint(favId), req.Name)
	if err != nil {
		logger.Logger.Errorf("[UpdateFavoriteName] failed to add favorite, err: %s", err.Error())
		c.JSON(http.StatusInternalServerError,  weberror.Info{Error: http.StatusText(http.StatusInternalServerError)})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// 删除收藏夹
// @Tags favorite
// @Description
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string          true       "userToken"
// @Param        fav_id         path      uint            true       "fav_id"
// @Success      204
// @Failure      400            {object}  weberror.Info              "Bad Request"
// @Failure      404            {object}  weberror.Info              "Not Found"
// @Failure      500            {object}  weberror.Info              "InternalServerError"
// @Router       /users/favorites/{fav_id} [delete]
// @Security     Token
func DeleteFavorite(c *gin.Context) {
	favId, err := strconv.ParseUint(c.Param("fav_id"), 10, 32)
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

	err = handler.DeleteFavorite(uint(user.ID), uint(favId))
	if err != nil {
		logger.Logger.Errorf("[DeleteFavorite] failed to delete favorite, err: %s", err.Error())
		c.JSON(http.StatusInternalServerError,  weberror.Info{Error: http.StatusText(http.StatusInternalServerError)})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// 获取收藏夹信息
// @Tags favorite
// @Description
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string          true       "userToken"
// @Param        fav_id         path      uint            true       "fav_id"
// @Success      200            {object}  payloads.GetFavoritesResp
// @Failure      400            {object}  weberror.Info              "Bad Request"
// @Failure      404            {object}  weberror.Info              "Not Found"
// @Failure      500            {object}  weberror.Info              "InternalServerError"
// @Router       /users/favorites/{fav_id} [get]
// @Security     Token
func GetFavorite(c *gin.Context) {
	favId, err := strconv.ParseUint(c.Param("fav_id"), 10, 32)
	if err != nil {
		logger.Logger.Errorf("[GetFavorite] failed to parse request, err: %s", err.Error())
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	user := extractclaims.ToUser(jwt.ExtractClaims(c))
	if user == nil {
		logger.Logger.Errorf("[GetFavorite] failed to parse request, err: %s", "failed to extract user info")
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}

	favorite, err := handler.QueryFavorite(uint(favId), "")
	if err != nil {
		logger.Logger.Errorf("[GetFavorite] failed to query favorite, err: %s", err.Error())
		c.JSON(http.StatusInternalServerError,  weberror.Info{Error: http.StatusText(http.StatusInternalServerError)})
		return
	}
	if favorite.UserId != user.ID {
		logger.Logger.Errorf("[GetFavorite] err: %s", "permission denied")
		c.JSON(http.StatusForbidden,  weberror.Info{Error: http.StatusText(http.StatusForbidden)})
		return
	}
	resp := payloads.GetFavoriteResp{
		Name: favorite.Name,
	}
	c.JSON(http.StatusOK, resp)
}

// 获取全部收藏夹
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
		logger.Logger.Errorf("[DeleteFavorite] failed to parse request, err: %s", err.Error())
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	offset, err := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 32)
	if err != nil {
		logger.Logger.Errorf("[DeleteFavorite] failed to parse request, err: %s", err.Error())
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	if limit < 1 || offset < 0 {
		logger.Logger.Errorf("[DeleteFavorite] failed to parse request,err: %s", "Parameter is out of range")
		c.JSON(http.StatusBadRequest, weberror.Info{Error: "Parameter is out of range"})
		return
	}
	user := extractclaims.ToUser(jwt.ExtractClaims(c))
	if user == nil {
		logger.Logger.Errorf("[DeleteFavorite] failed to parse request, err: %s", "failed to extract user info")
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}

	favorites, total, err := handler.QueryFavorites(user.ID, uint(limit), uint(offset))
	if err != nil {
		logger.Logger.Errorf("[DeleteFavorite] failed to query favorites, err: %s", err.Error())
		c.JSON(http.StatusInternalServerError,  weberror.Info{Error: http.StatusText(http.StatusInternalServerError)})
		return
	}
	resps := payloads.GetFavoritesResp{
		Total: total,
		Favs: make([]payloads.Favorite, 0, len(favorites)),
	}
	for _, v := range favorites {
		resps.Favs = append(resps.Favs, payloads.Favorite{
			FavId: v.ID,
			Name: v.Name,
		})
	}
	c.JSON(http.StatusOK, resps)
}

// 添加收藏
// @Tags favorite
// @Description
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string              true  "userToken"
// @Param        fav_id         path      uint                true  "fav_id"
// @Param        AddDocReq      body      payloads.AddDocReq  true  "include doc_index"
// @Success      201
// @Failure      400            {object}  weberror.Info             "Bad Request"
// @Failure      404            {object}  weberror.Info             "Not Found"
// @Failure      500            {object}  weberror.Info             "InternalServerError"
// @Router       /users/favorites/{fav_id}/docs [put]
// @Security     Token
func AddDoc(c *gin.Context) {
	favId, err := strconv.ParseUint(c.Param("fav_id"), 10, 32)
	if err != nil {
		logger.Logger.Errorf("[AddDoc] failed to parse request, err: %s", err.Error())
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	var req payloads.AddDocReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Logger.Errorf("[AddDoc] failed to parse request, err: %s", err.Error())
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	user := extractclaims.ToUser(jwt.ExtractClaims(c))
	if user == nil {
		logger.Logger.Errorf("[AddDoc] failed to parse request, err: %s", "failed to extract user info")
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}

	//是否有对此收藏夹的访问权限
	fav, err := handler.QueryFavorite(uint(favId), "")
	if err != nil {
		logger.Logger.Errorf("[AddDoc] failed to query favorite, err: %s", err.Error())
		c.JSON(http.StatusInternalServerError,  weberror.Info{Error: http.StatusText(http.StatusInternalServerError)})
		return
	}
	if fav.UserId != user.ID {
		logger.Logger.Errorf("[AddDoc] err: %s", "permission denied")
		c.JSON(http.StatusForbidden,  weberror.Info{Error: http.StatusText(http.StatusForbidden)})
		return
	}

	//此文档是否已收藏
	_, err = handler.QueryDoc(0, req.DocIndex)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		if err != nil {
			logger.Logger.Errorf("[AddDoc] failed to query doc, err: %s", err.Error())
			c.JSON(http.StatusInternalServerError,  weberror.Info{Error: http.StatusText(http.StatusInternalServerError)})
			return
		}
		logger.Logger.Errorf("[AddDoc] failed to append doc, err: %s", "duplicate document")
		c.JSON(http.StatusInternalServerError, weberror.Info{Error: "duplicate document"})
		return
	}

	err = handler.AppendDoc(uint(favId), req.DocIndex)
	if err != nil {
		logger.Logger.Errorf("[AddDoc] failed to append doc, err: %s", err.Error())
		c.JSON(http.StatusInternalServerError,  weberror.Info{Error: http.StatusText(http.StatusInternalServerError)})
		return
	}
	c.JSON(http.StatusCreated, nil)
}

// 取消收藏
// @Tags favorite
// @Description
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string              true   "userToken"
// @Param        fav_id         path      uint                true   "fav_id"
// @Param        doc_id         path      uint                true   "doc_id"
// @Success      204
// @Failure      400            {object}  weberror.Info              "Bad Request"
// @Failure      404            {object}  weberror.Info              "Not Found"
// @Failure      500            {object}  weberror.Info              "InternalServerError"
// @Router       /users/favorites/{fav_id}/docs/{doc_id} [delete]
// @Security     Token
func DeleteDoc(c *gin.Context) {
	favId, err := strconv.ParseUint(c.Param("fav_id"), 10, 32)
	if err != nil {
		logger.Logger.Errorf("[DeleteDoc] failed to parse request, err: %s", err.Error())
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	docId, err := strconv.ParseUint(c.Param("doc_id"), 10, 32)
	if err != nil {
		logger.Logger.Errorf("[DeleteDoc] failed to parse request, err: %s", err.Error())
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	user := extractclaims.ToUser(jwt.ExtractClaims(c))
	if user == nil {
		logger.Logger.Errorf("[DeleteDoc] failed to parse request, err: %s", "failed to extract user info")
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}

	//是否有对此收藏夹的访问权限
	fav, err := handler.QueryFavorite(uint(favId), "")
	if err != nil {
		logger.Logger.Errorf("[DeleteDoc] failed to query favorite, err: %s", err.Error())
		c.JSON(http.StatusInternalServerError,  weberror.Info{Error: http.StatusText(http.StatusInternalServerError)})
		return
	}
	if fav.UserId != user.ID {
		logger.Logger.Errorf("[DeleteDoc] err: %s", "permission denied")
		c.JSON(http.StatusForbidden,  weberror.Info{Error: http.StatusText(http.StatusForbidden)})
		return
	}

	err = handler.DeleteDoc(uint(favId), uint(docId))
	if err != nil {
		logger.Logger.Errorf("[DeleteDoc] failed to delete doc, err: %s", err.Error())
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
// @Param        fav_id         path      uint            true       "fav_id"
// @Param        limit          query     uint            false      "greater than 0"
// @Param        offset         query     uint            false		 "greater than -1"
// @Success      200            {object}  payloads.GetDocsResp
// @Failure      400            {object}  weberror.Info              "Bad Request"
// @Failure      404            {object}  weberror.Info              "Not Found"
// @Failure      500            {object}  weberror.Info              "InternalServerError"
// @Router       /users/favorites/{fav_id}/docs [get]
// @Security     Token
func GetDocs(c *gin.Context) {
	favId, err := strconv.ParseUint(c.Param("fav_id"), 10, 32)
	if err != nil {
		logger.Logger.Errorf("[GetDocs] failed to parse request, err: %s", err.Error())
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
	}
	limit, err := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 32)
	if err != nil {
		logger.Logger.Errorf("[GetDocs] failed to parse request, err: %s", err.Error())
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	offset, err := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 32)
	if err != nil {
		logger.Logger.Errorf("[GetDocs] failed to parse request, err: %s", err.Error())
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	if limit < 1 || offset < 0 {
		logger.Logger.Errorf("[GetDocs] failed to parse request,err: %s", "Parameter is out of range")
		c.JSON(http.StatusBadRequest, weberror.Info{Error: "Parameter is out of range"})
		return
	}
	user := extractclaims.ToUser(jwt.ExtractClaims(c))
	if user == nil {
		logger.Logger.Errorf("[GetDocs] failed to parse request, err: %s", "failed to extract user info")
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}

	//是否有对此收藏夹的访问权限
	fav, err := handler.QueryFavorite(uint(favId), "")
	if err != nil {
		logger.Logger.Errorf("[GetDocs] failed to query favorite, err: %s", err.Error())
		c.JSON(http.StatusInternalServerError,  weberror.Info{Error: http.StatusText(http.StatusInternalServerError)})
		return
	}
	if fav.UserId != user.ID {
		logger.Logger.Errorf("[GetDocs] err: %s", "permission denied")
		c.JSON(http.StatusForbidden,  weberror.Info{Error: http.StatusText(http.StatusForbidden)})
		return
	}

	docs, total, err := handler.QueryDocs(uint(favId), uint(limit), uint(offset))
	if err != nil {
		logger.Logger.Errorf("[GetDocs] failed to query docs, err: %s", err.Error())
		c.JSON(http.StatusInternalServerError,  weberror.Info{Error: http.StatusText(http.StatusInternalServerError)})
		return
	}
	resps := payloads.GetDocsResp{
		Total: total,
		Docs: make([]payloads.Doc, 0, len(docs)),
	}
	for _, v := range docs {
		resps.Docs = append(resps.Docs, payloads.Doc{
			DocId: v.ID,
			DocIndex: v.DocIndex,
			Url: v.Url,
			Summary: v.Summary,
		})
	}
	c.JSON(http.StatusOK, resps)
}