package model

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type Principal struct {
	UserID   uint   `json:"user_id"`
	UserName string `json:"username"`
	UserType string `json:"user_type"` //"user"和"admin""
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

// 真为管理，假为用户
func (p *Principal) IsAdmin() bool {
	if p.UserType == "admin" {
		return true
	}
	return false
}
