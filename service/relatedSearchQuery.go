package service

import (
	"lgdSearch/pkg/trie"
)

type PostInfo struct {
	Text string
}
type PageInfo struct {
	PostList []*PostInfo
}

func QueryPageInfo(text string) (*PageInfo, error) {
	postList := getPostList(text)
	return &PageInfo{
		postList,
	}, nil
}

func getPostList(text string) []*PostInfo {
	postList := make([]*PostInfo, 0)
	result := trie.Tree.Search([]rune(text))
	for _, val := range result {
		postList = append(postList, &PostInfo{
			Text: val,
		})
	}
	return postList
}
