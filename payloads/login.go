package payloads

import (
	"time"
)

type LoginReq struct {
	Username 	string `form:"username" json:"username" binding:"required"`
	Password 	string `form:"password" json:"password" binding:"required"`
}

type LoginResp struct {
	Expire 		time.Time `form:"expire" json:"expire" binding:"required"`
	Token 		string `form:"token" json:"token" binding:"required"`
}