package payloads

type AddDocReq struct {
	DocIndex  uint `form:"doc_index" json:"doc_index" binging:"required"`
}