package model

import "errors"

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

// 预定义的错误变量
var (
	ErrUserNotFound      = errors.New("用户不存在")
	ErrUserAlreadyExists = errors.New("用户已存在")
	ErrPasswordWrong     = errors.New("密码错误")
	ErrBookNotFound      = errors.New("图书不存在")
	ErrInsufficientStock = errors.New("图书库存不足")
	ErrNoUnreturnedBooks = errors.New("用户有未归还的图书")
)
