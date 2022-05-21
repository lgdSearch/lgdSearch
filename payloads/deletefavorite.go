package payloads

type DeleteFavoriteReq struct {
	DocId 	uint `form:"doc_id" json:"doc_id" binging:"required"`
}