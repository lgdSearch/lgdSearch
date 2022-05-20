package payloads

type UpdateProfileReq struct {
	Nickname 	string `form:"nickname" json:"nickname" binging:"required"`
}