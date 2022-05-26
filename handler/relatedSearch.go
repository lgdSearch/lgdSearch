package handler

import (
	"lgdSearch/pkg/trie"
)

type PageData struct {
	Code int64       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func RelatedQuery(text string) (*PageData, error) {
	textInfo, err := QueryPageInfo(text)
	if err != nil {
		return &PageData{
			Code: -1,
			Msg:  err.Error(),
		}, err
	}
	return &PageData{
		Code: 0,
		Msg:  "success",
		Data: textInfo,
	}, nil
}

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
