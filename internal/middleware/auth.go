package middleware

import (
	"strings"

	"github.com/Cheerdoge/library-manage-system/internal/model"
	"github.com/Cheerdoge/library-manage-system/internal/service"
	"github.com/Cheerdoge/library-manage-system/web"
	"github.com/gin-gonic/gin"
)

var jwtService = service.NewJWTService()

func AuthMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(model.CodeUnauthorized, web.ErrorResponse(model.CodeUnauthorized, "没有令牌"))
		return
	}

	// 解析 Bearer 前缀
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.AbortWithStatusJSON(model.CodeUnauthorized, web.ErrorResponse(model.CodeUnauthorized, "Authorization 格式错误，应为 Bearer <token>"))
		return
	}

	token := parts[1]
	principal, err := jwtService.ParseToken(token)
	if err != nil {
		c.AbortWithStatusJSON(model.CodeUnauthorized, web.ErrorResponse(model.CodeUnauthorized, "令牌无效: "+err.Error()))
		return
	}
	c.Set("principal", principal) //注入上下文
	c.Next()
}
