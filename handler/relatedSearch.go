package handler

import "lgdSearch/service"

type PageData struct {
	Code int64       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func RelatedQuery(text string) (*PageData, error) {
	textInfo, err := service.QueryPageInfo(text)
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
