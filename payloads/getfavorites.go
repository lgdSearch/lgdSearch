package payloads

type GetFavoritesResp struct {
	FavId 	uint `form:"doc_id" json:"doc_id"`
	Name    string `form:"summary" json:"summary"`
}