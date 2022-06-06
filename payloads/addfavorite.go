package payloads

type AddFavoriteReq struct {
	Name  string `form:"name" json:"name" binging:"required"`
}