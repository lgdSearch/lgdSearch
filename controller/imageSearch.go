package controller

import (
	"lgdSearch/payloads"
	"lgdSearch/pkg/logger"
	"lgdSearch/pkg/vgg"
	"lgdSearch/pkg/weberror"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
    "image/jpeg"
	"bytes"
)

// 以图搜图
// @Tags search
// @Description
// @Accept       json
// @Produce      json
// @Success      200            {object}  payloads.ImageSearchResp
// @Failure      400            {object}  weberror.Info              "Bad Request"
// @Failure      404            {object}  weberror.Info              "Not Found"
// @Failure      500            {object}  weberror.Info              "InternalServerError"
// @Router       /image_search [post]
// @Security     Token
func ImageSearch(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		logger.Logger.Errorf("[ImageSearch] failed to parse request, err: %s", err.Error())
		c.JSON(http.StatusBadRequest, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	defer file.Close()
	// img, err := ioutil.ReadAll(file)
	// if err != nil {
	// 	logger.Logger.Errorf("[ImageSearch] failed to read file, err: %s", err.Error())
	// 	c.JSON(http.StatusInternalServerError, weberror.Info{Error: http.StatusText(http.StatusInternalServerError)})
	// 	return
	// }
	img, err := jpeg.Decode(file)
	if err != nil {
		logger.Logger.Errorf("[ImageSearch] failed to decode, err: %s", err.Error())
		c.JSON(http.StatusInternalServerError, weberror.Info{Error: http.StatusText(http.StatusInternalServerError)})
		return
	}
	newImage := resize.Resize(224, 224, img, resize.NearestNeighbor)
	if err != nil {
		logger.Logger.Errorf("[ImageSearch] failed to encode, err: %s", err.Error())
		c.JSON(http.StatusInternalServerError, weberror.Info{Error: http.StatusText(http.StatusInternalServerError)})
		return
	}
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, newImage, nil)
	if err != nil {
		logger.Logger.Errorf("[ImageSearch] failed to encode, err: %s", err.Error())
		c.JSON(http.StatusInternalServerError, weberror.Info{Error: http.StatusText(http.StatusInternalServerError)})
		return
	}
	bImage := buf.Bytes()
	imgs, err := vgg.Search(bImage)
	if err != nil {
		logger.Logger.Errorf("[ImageSearch] failed to search, err: %s", err.Error())
		c.JSON(http.StatusInternalServerError, weberror.Info{Error: http.StatusText(http.StatusInternalServerError)})
		return
	}
	resp := &payloads.ImageSearchResp{
		Images: imgs,
	}
	c.JSON(http.StatusOK, resp)
}