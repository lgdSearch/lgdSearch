package handler

import (
	"lgdSearch/pkg"
	"lgdSearch/pkg/models"
	"lgdSearch/pkg/trie"
)

func MultiSearch(request *models.SearchRequest) *models.SearchResult {
	trie.Tree.InsertString(request.Query)
	trie.SendText(request.Query)

	return pkg.SearchEngine.MultiSearch(request)
}
