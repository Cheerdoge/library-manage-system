package model

// 错误码定义 (HTTP状态码)
const (
	// 成功
	CodeSuccess = 200

	// 4xx 客户端错误
	CodeBadRequest   = 400
	CodeUnauthorized = 401
	CodeForbidden    = 403
	CodeNotFound     = 404
	CodeConflict     = 409

	// 5xx 服务器错误
	CodeServerError = 500
)

const (
	ErrUserAlreadyExists = "用户已存在"
	ErrUserNotFound      = "用户不存在"
	ErrPasswordWrong     = "密码错误"
	ErrInvalidInput      = "无效的输入参数"
	ErrUnauthorized      = "未授权的操作"
	ErrForbidden         = "无权限执行此操作"
	ErrConflict          = "请求与当前资源状态冲突"
	ErrServerInternal    = "服务器内部错误"
)
