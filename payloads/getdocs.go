package payloads

type GetDocsResp struct {
	Total   int64   `form:"total" json:"total"`
	Docs    []Doc `form:"docs" json:"docs"`
}

type Doc struct {
	DocId    uint   `form:"doc_id" json:"doc_id"`
	DocIndex uint   `form:"doc_index" json:"doc_index"`
	Url      string `form:"url" json:"url"` 
	Summary  string `form:"summary" json:"summary"`
}