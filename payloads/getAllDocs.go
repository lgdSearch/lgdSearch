package payloads

type GetAllDocsResp struct {
	Favs    []FavoriteWithDocs `form:"favs" json:"favs"`
}

type FavoriteWithDocs struct {
	Favorite
	Docs 	[]Doc `form:"docs" json:"docs"`
}