package handler

import (
	"strconv"

	"github.com/Cheerdoge/library-manage-system/internal/middleware"
	"github.com/Cheerdoge/library-manage-system/internal/model"
	"github.com/Cheerdoge/library-manage-system/internal/service"
	"github.com/Cheerdoge/library-manage-system/web"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	sessionservice *service.SessionService
	userservice    *service.UserService
}

func NewAuthHandler(sessionservice *service.SessionService, userservice *service.UserService) *AuthHandler {
	return &AuthHandler{
		sessionservice: sessionservice,
		userservice:    userservice,
	}
}

// LoginHandler 登录
func (h *AuthHandler) LoginHandler(c *gin.Context) {
	var req web.LoginUser
	err := c.ShouldBindJSON(&req)
	if err != nil {
		web.FailWithMessage(c, "请求参数有误")
		return
	}

	if req.Username == "" || req.Password == "" {
		web.FailWithMessage(c, "用户名或密码不能为空")
		return
	}

	token, msg := h.userservice.Login(req.Username, req.Password)
	if msg != "" {
		web.FailWithMessage(c, msg)
		return
	}

	c.SetCookie("session_token", token, int(middleware.SessionDuration.Seconds()), "/", "", false, true)
	web.OkWithMessage(c, "登录成功")
}

// LogoutHandler 登出
func (h *AuthHandler) LogoutHandler(c *gin.Context) {
	token, err := c.Cookie("session_token")
	if err != nil {
		web.FailWithMessage(c, "获取cookie中session失败")
		return
	}
	msg := h.sessionservice.DelSessionByToken(token)
	if msg != "" {
		web.FailWithMessage(c, "删除session失败:"+msg)
		return
	}
	web.OkWithMessage(c, "登出成功")
}

// DelHandler 注销用户
func (h *AuthHandler) DelHandler(c *gin.Context) {
	var req web.DelUser
	err := c.ShouldBindJSON(&req)
	if err != nil {
		web.FailWithMessage(c, "请求参数有误")
		return
	}
	if req.Username == "" || req.Password == "" {
		web.FailWithMessage(c, "用户名或密码不能为空")
		return
	}
	token, err := c.Cookie("session_token")
	if err != nil {
		web.FailWithMessage(c, "获取cookie中session失败")
		return
	}
	_, msg := h.sessionservice.CheckSessionByToken(token)
	if msg != "" {
		web.FailWithMessage(c, msg)
		return
	}

	principal, err := model.GetPrincipal(c)
	if err != nil {
		web.FailWithMessage(c, "无法获取用户信息")
		return
	}
	targetUsername := req.Username
	if !principal.IsAdmin && targetUsername != principal.UserName {
		web.FailWithMessage(c, "无权注销他人账户")
		return
	}

	msg = h.userservice.WithdrawUser(targetUsername, req.Password)
	if msg != "" {
		web.FailWithMessage(c, msg)
		return
	}

	web.OkWithMessage(c, "注销成功")
}

// RegisterUserHandler 注册新用户
func (h *AuthHandler) RegisterUserHandler(c *gin.Context) {
	var req web.AddUser
	err := c.ShouldBindJSON(&req)
	if err != nil {
		web.FailWithMessage(c, "请求参数有误")
		return
	}
	if req.Username == "" || req.Password == "" {
		web.FailWithMessage(c, "用户名或密码不能为空")
		return
	}
	id, msg := h.userservice.Register(req.Username, req.Password, false)
	if msg != "" {
		web.FailWithMessage(c, msg)
		return
	}
	web.Ok(c, "注册成功,请牢记账户和密码", gin.H{"user_id": id})
}

// RegisterAdminHandler 注册管理员用户
func (h *AuthHandler) RegisterAdminHandler(c *gin.Context) {
	var req web.AddUser
	err := c.ShouldBindJSON(&req)
	if err != nil {
		web.FailWithMessage(c, "请求参数有误")
		return
	}
	id, msg := h.userservice.Register(req.Username, req.Password, true)
	if msg != "" {
		web.FailWithMessage(c, msg)
		return
	}
	web.Ok(c, "注册成功,请牢记账户和密码", gin.H{"user_id": id})
}

// 管理员注销用户 AdminDelHandler
func (h *AuthHandler) AdminDelHandler(c *gin.Context) {
	useridstr := c.Param("userid")
	userid, err := strconv.Atoi(useridstr)
	if err != nil {
		web.FailWithMessage(c, "请求参数有误")
		return
	}
	msg := h.userservice.AdminWithdrawUser(uint(userid))
	if msg != "" {
		web.FailWithMessage(c, msg)
		return
	}
	web.OkWithMessage(c, "用户注销成功")
}
