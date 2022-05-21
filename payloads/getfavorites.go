package payloads

type GetFavoritesResp struct {
	DocId 	uint `form:"doc_id" json:"doc_id"`
	Summary string `form:"summary" json:"summary"`
}