package payloads

type LoginResq struct {
	Username 	string `form:"username" json:"username" binding:"required"`
	Password 	string `form:"password" json:"password" binding:"password"`
}