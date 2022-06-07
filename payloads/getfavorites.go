package payloads

type GetFavoritesResp struct {
	Total  int64   `form:"total" json:"total"`
	Favs   []Favorite `form:"favs" json:"favs"`
}

type Favorite struct {
	FavId 	uint `form:"fav_id" json:"fav_id"`
	Name    string `form:"name" json:"name"`
}