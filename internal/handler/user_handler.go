package handler

import (
	"github.com/Cheerdoge/library-manage-system/internal/model"
	"github.com/Cheerdoge/library-manage-system/internal/service"
	"github.com/Cheerdoge/library-manage-system/web"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userservice *service.UserService
}

func NewUserHandler(userservice *service.UserService) *UserHandler {
	return &UserHandler{
		userservice: userservice,
	}
}

// GetUserInfoHandler 获取用户信息
func (h *UserHandler) GetUserInfoHandler(c *gin.Context) {
	principal, err := model.GetPrincipal(c)
	if err != nil {
		web.FailWithMessage(c, "无法从上下文获取用户信息")
		return
	}
	userinfo, msg := h.userservice.GetUserInfo(principal.UserID)
	if msg != "" {
		web.FailWithMessage(c, msg)
		return
	}
	web.OkWithData(c, userinfo)
}

// UserChangePasswordHandler 用户修改用户密码
func (h *UserHandler) UserChangePasswordHandler(c *gin.Context) {
	var req web.ChangePassword
	err := c.ShouldBindJSON(&req)
	if err != nil {
		web.FailWithMessage(c, "请求参数有误")
		return
	}
	principal, err := model.GetPrincipal(c)
	if err != nil {
		web.FailWithMessage(c, "无法获取用户信息")
		return
	}

	msg := h.userservice.ChangePassword(principal.IsAdmin, principal.UserID, req.Oldpassword, req.Newpassword)
	if msg != "" {
		web.FailWithMessage(c, msg)
		return
	}

	web.OkWithMessage(c, "密码修改成功")
}

// AdminChangePasswordHandler 管理员重置用户密码
func (h *UserHandler) AdminChangePasswordHandler(c *gin.Context) {
	var req web.ChangePassword
	err := c.ShouldBindJSON(&req)
	if err != nil {
		web.FailWithMessage(c, "请求参数有误")
		return
	}
	principal, err := model.GetPrincipal(c)
	if err != nil {
		web.FailWithMessage(c, "无法获取用户信息")
		return
	}
	msg := h.userservice.ChangePassword(principal.IsAdmin, req.UserId, "", req.Newpassword)
	if msg != "" {
		web.FailWithMessage(c, msg)
		return
	}
	web.OkWithMessage(c, "密码重置成功")
}

// ChangeUserInfoHandler 修改用户信息
func (h *UserHandler) ChangeUserInfoHandler(c *gin.Context) {
	var req web.ChangeUserInfo
	err := c.ShouldBindJSON(&req)
	if err != nil {
		web.FailWithMessage(c, "请求参数有误")
		return
	}
	principal, err := model.GetPrincipal(c)
	if err != nil {
		web.FailWithMessage(c, "无法通过上下文获取用户信息")
		return
	}
	msg := h.userservice.ChangeUserInfo(principal.UserID, req.Username, req.Telenum)
	if msg != "" {
		web.FailWithMessage(c, msg)
		return
	}
	web.OkWithMessage(c, "用户信息修改成功")
}

// AdminGetUserInfoHandler 管理员获取指定用户信息
func (h *UserHandler) AdminGetUserInfoHandler(c *gin.Context) {
	var req web.GetUserInfo
	err := c.ShouldBindJSON(&req)
	if err != nil {
		web.FailWithMessage(c, "请求参数有误")
		return
	}
	principal, err := model.GetPrincipal(c)
	if err != nil {
		web.FailWithMessage(c, "无法通过上下文获取用户信息")
		return
	}
	if !principal.IsAdmin {
		web.FailWithMessage(c, "权限不足")
		return
	}
	userinfo, msg := h.userservice.GetUserInfoByName(req.Username)
	if msg != "" {
		web.FailWithMessage(c, msg)
		return
	}
	web.OkWithData(c, userinfo)
}

// AdminGetAllUserInfoHandler 管理员获取所有用户信息
func (h *UserHandler) AdminGetAllUserInfoHandler(c *gin.Context) {
	principal, err := model.GetPrincipal(c)
	if err != nil {
		web.FailWithMessage(c, "无法通过上下文获取用户信息")
		return
	}
	userinfos, msg := h.userservice.GetAllUsersInfo(principal.IsAdmin)
	if msg != "" {
		web.FailWithMessage(c, msg)
		return
	}
	web.OkWithData(c, userinfos)
}
