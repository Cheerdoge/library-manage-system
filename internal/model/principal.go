package model

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type Principal struct {
	UserID   uint   `json:"user_id"`
	UserName string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
}

func GetPrincipal(c *gin.Context) (*Principal, error) {
	principal, exists := c.Get("principal")
	if !exists {
		return nil, errors.New("无法获取用户信息")
	}
	p, ok := principal.(*Principal)
	if !ok {
		return nil, errors.New("用户信息类型错误")
	}
	return p, nil
}
