package entities

type BaseResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

func ResSuccess[T any](data T) BaseResponse {
	return BaseResponse{
		Code: 0,
		Data: data,
	}
}

func ResFailed(code int, msg string) BaseResponse {
	return BaseResponse{
		Code: code,
		Msg:  msg,
	}
}
