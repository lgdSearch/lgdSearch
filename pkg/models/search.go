package models

// Highlight 关键词高亮
type Highlight struct {
	PreTag  string `json:"preTag"`  //高亮前缀
	PostTag string `json:"postTag"` //高亮后缀
}

// SearchRequest 搜索请求
type SearchRequest struct {
	Query      string     `json:"query,omitempty"`      // 搜索关键词
	Page       int        `json:"page,omitempty"`       // 页码
	Limit      int        `json:"limit,omitempty"`      // 每页大小，最大1000，超过报错
	FilterWord []string   `json:"filterWord,omitempty"` // 关键词过滤
	Highlight  *Highlight `json:"highlight,omitempty"`  // 关键词高了
}

func (s *SearchRequest) GetAndSetDefault() *SearchRequest {

	if s.Limit == 0 {
		s.Limit = 20
	}
	if s.Page == 0 {
		s.Page = 1
	}

	if s.Highlight == nil {
		s.Highlight = &Highlight{
			PreTag:  "<span style=\"color: red;\">",
			PostTag: "</span>",
		}
	}

	return s
}

// SearchResult 搜索响应
type SearchResult struct {
	Time      float32       `json:"time,omitempty"`      //查询用时
	Total     int           `json:"total"`               //总数
	PageCount int           `json:"pageCount"`           //总页数
	Page      int           `json:"page,omitempty"`      //页码
	Limit     int           `json:"limit,omitempty"`     //页大小
	Documents []ResponseDoc `json:"documents,omitempty"` //文档
	Related   []string      `json:"related,omitempty"`   // 相关搜索
	Words     []string      `json:"words,omitempty"`     //搜索关键词
}

type SearchPictureResult struct {
	Time      float32       `json:"time,omitempty"`      //查询用时
	Total     int           `json:"total"`               //总数
	PageCount int           `json:"pageCount"`           //总页数
	Page      int           `json:"page,omitempty"`      //页码
	Limit     int           `json:"limit,omitempty"`     //页大小
	Documents []ResponseUrl `json:"documents,omitempty"` //缩略图 Url
	Words     []string      `json:"words,omitempty"`     //搜索关键词
}
