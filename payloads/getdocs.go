package payloads

type GetDocsResp struct {
	DocId    uint   `form:"doc_id" json:"doc_id"`
	DocIndex uint   `form:"doc_index" json:"doc_index"`
	Summary  string `form:"summary" json:"summary"`
}