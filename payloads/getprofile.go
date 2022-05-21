package payloads

type GetProfileResp struct {
	Username 	string `form:"username" json:"username"`
	Nickname 	string `form:"nickname" json:"nickname"`
}