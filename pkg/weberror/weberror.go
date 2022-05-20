package weberror

type Info struct {
	Code  int    `json:"code"`  //错误码，用于前端依靠http状态码不能区别错误信息时用
	Error string `json:"error"` //系统错误详细信息
}