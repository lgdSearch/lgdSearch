package payloads

type GetFavoriteResp struct {
	Name string `form:"name" json:"name"`
}