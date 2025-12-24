package service

import (
	"github.com/Cheerdoge/library-manage-system/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	AddUser(username string, password string, isadmin bool) (uint, error)
	FindUserById(Id uint) (*model.User, error)
	UpdatePassword(userid uint, newpassword string) error
	UpdateUserInfo(userid uint, username string, telenum string, overdueNum int) error
	FindUserByName(name string) (*model.User, error)
	DeleUser(userid uint) error
	GetAllUsers() ([]model.UserInfo, error)
	ModifyUserNum(tx *gorm.DB, userid uint, borrowChange int, overdueChange int) error
}

type UserService struct {
	userrepo       UserRepository
	sessionservice *SessionService
}

func NewUserService(repo UserRepository, sessionservice *SessionService) *UserService {
	return &UserService{
		userrepo:       repo,
		sessionservice: sessionservice,
	}
}

// Register 注册新用户
// 如果id为0，则表示注册失败
func (s *UserService) Register(username string, password string, isadmin bool) (id uint, message string) {
	_, err := s.userrepo.FindUserByName(username)
	if err == nil {
		return 0, model.ErrUserAlreadyExists
	}
	id, err = s.userrepo.AddUser(username, password, isadmin)
	if err != nil {
		return 0, model.ErrServerInternal
	}
	return id, ""
}

// Login 用户登录
// token为空即失败
func (s *UserService) Login(username string, password string) (token string, message string) {
	targetuser, err := s.userrepo.FindUserByName(username)
	if err != nil {
		return "", model.ErrUserNotFound
	}
	if targetuser.Password != password {
		return "", model.ErrPasswordWrong
	}

	token, msg := s.sessionservice.CreateSession(targetuser.ID, targetuser.UserName, targetuser.IsAdmin)
	if msg != "" {
		return "", msg
	}

	return token, ""
}

// Logout 用户登出
// 成功返回空字符串
func (s *UserService) Logout(userid uint) (message string) {
	session, err := s.sessionservice.repo.FindSessionByUserId(userid)
	if err != nil {
		return model.ErrUserNotFound
	}
	err = s.sessionservice.repo.DeleteSessionByToken(session.Token)
	if err != nil {
		return "登录已过期"
	}
	return ""
}

// ChangePassword 修改密码
// 管理员修改用户密码，无需验证旧密码
// 成功返回空字符串
func (s *UserService) ChangePassword(isadmin bool, userid uint, oldpassword string, newpassword string) (message string) {
	if !isadmin {
		user, err := s.userrepo.FindUserById(userid)
		if err != nil {
			return model.ErrUserNotFound
		}
		if user.Password != oldpassword {
			return model.ErrPasswordWrong
		}
	}
	err := s.userrepo.UpdatePassword(userid, newpassword)
	if err != nil {
		return err.Error()
	}
	return ""
}

// GetUserInfo 获取用户信息
// 失败返回nil和错误信息
func (s *UserService) GetUserInfo(userid uint) (user *model.UserInfo, message string) {
	targetuser, err := s.userrepo.FindUserById(userid)
	if err != nil {
		return nil, model.ErrUserNotFound
	}
	return &model.UserInfo{
		ID:       targetuser.ID,
		UserName: targetuser.UserName,
		Telenum:  targetuser.Telenum,
		IsAdmin:  targetuser.IsAdmin,
	}, ""
}

// ChangeUserInfo 修改用户信息
// 成功返回空字符串
func (s *UserService) ChangeUserInfo(userid uint, username string, telenum string) (message string) {
	err := s.userrepo.UpdateUserInfo(userid, username, telenum, 0)
	if err != nil {
		return err.Error()
	}
	return ""
}

// WithdrawUser 删除用户
// 用户验证密码，管理员直接操作
// 成功返回空字符串
func (s *UserService) WithdrawUser(isadmin bool, username string, password string) (message string) {
	user, err := s.userrepo.FindUserByName(username)
	if err != nil {
		return model.ErrUserNotFound
	}
	if !isadmin {
		if username != user.UserName {
			return model.ErrForbidden
		}
		if user.Password != password {
			return model.ErrPasswordWrong
		}
	}
	if user.NowBorrNum > 0 {
		return "用户有未归还的书籍，无法删除"
	}
	err = s.userrepo.DeleUser(user.ID)
	if err != nil {
		return err.Error()
	}
	return ""
}

// GetAllUsersInfo 获取所有用户信息
// 管理员专属
// 失败返回空切片
func (s *UserService) GetAllUsersInfo(isadmin bool) (userlist []model.UserInfo, message string) {
	if !isadmin {
		return nil, model.ErrForbidden
	}
	userlist, err := s.userrepo.GetAllUsers()
	if err != nil {
		return nil, err.Error()
	}
	return userlist, ""
}

func (s *UserService) GetUserInfoByName(username string) (user *model.UserInfo, message string) {
	targetuser, err := s.userrepo.FindUserByName(username)
	if err != nil {
		return nil, model.ErrUserNotFound
	}
	return &model.UserInfo{
		ID:       targetuser.ID,
		UserName: targetuser.UserName,
		Telenum:  targetuser.Telenum,
		IsAdmin:  targetuser.IsAdmin,
	}, ""
}
