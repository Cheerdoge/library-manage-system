package web

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// 错误响应函数
func ErrorResponse(code int, message string) Response {
	return Response{
		Code:    code,
		Message: message,
		Data:    nil,
	}
}

// 成功响应函数
func SuccessResponse(message string, data interface{}) Response {
	return Response{
		Code:    200,
		Message: "成功",
		Data:    data,
	}
}
