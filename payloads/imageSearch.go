package payloads

type ImageSearchResp struct {
	Images    [][]byte `form:"images" json:"images"`
}