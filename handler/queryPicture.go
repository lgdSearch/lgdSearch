package handler

import (
	"lgdSearch/pkg"
	"lgdSearch/pkg/models"
)

func MultiSearchPicture(request *models.SearchRequest) *models.SearchPictureResult {
	return pkg.SearchEngine.MultiSearchPicture(request)
}
