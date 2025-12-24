package middleware

import (
	"time"

	"github.com/Cheerdoge/library-manage-system/internal/model"
	"github.com/Cheerdoge/library-manage-system/internal/service"
	"github.com/Cheerdoge/library-manage-system/web"
	"github.com/gin-gonic/gin"
)

const SessionDuration = 30 * time.Minute // Session有效期为7天

func AuthMiddleware(sessionservice *service.SessionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("session_token")
		if err != nil {
			web.Fail(c, model.CodeUnauthorized, "请先登录")
			c.Abort()
			return
		}

		session, msg := sessionservice.CheckSessionByToken(token)
		if msg != "" {
			web.Fail(c, model.CodeUnauthorized, msg)
			c.Abort()
			return
		}

		principal := &model.Principal{
			UserID:   session.UserID,
			UserName: session.UserName,
			IsAdmin:  session.IsAdmin,
		}
		c.Set("principal", principal)
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		principal, err := model.GetPrincipal(c)
		if err != nil {
			web.Fail(c, model.CodeUnauthorized, "无法获取用户信息")
			c.Abort()
			return
		}
		if !principal.IsAdmin {
			web.Fail(c, model.CodeForbidden, "权限不足")
			c.Abort()
			return
		}
		c.Next()
	}
}
