package payloads

type AddFavoriteReq struct {
	DocId 	uint `form:"doc_id" json:"doc_id" binging:"required"`
}