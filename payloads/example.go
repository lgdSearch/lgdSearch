package payloads

type SayHelloReq struct {
	Text 	string `json:"text"`
}

type SayHelloResp struct {
	Text 	string `json:"text"`
}