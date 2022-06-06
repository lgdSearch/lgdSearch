package payloads

type GetFavoriteResp struct {
	Name string `form:"name" json:"name"`
	Docs []Doc  `form:"docs" json:"docs"`
}


type Doc struct {
	DocId   uint
	Summary string
}