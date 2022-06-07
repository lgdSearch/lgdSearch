package payloads

type UpdateFavoriteNameReq struct {
	Name string `form:"name" json:"name" binding:"required"`
}