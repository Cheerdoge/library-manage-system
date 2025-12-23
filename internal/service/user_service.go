package service

import (
	"github.com/Cheerdoge/library-manage-system/internal/model"
)

type UserRepository interface {
	AddUser(username string, password string, usertype string) (uint, error)
	FindUserById(Id uint) (*model.User, error)
	UpdatePassword(userid uint, newpassword string) error
	UpdateUserInfo(userid uint, username string, telenum string) error
	FindUserByName(name string) (*model.User, error)
	DeleUser(userid uint) error
	GetAllUsers() ([]model.UserInfo, error)
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// Register 注册新用户
// 如果id为0，则表示注册失败
func (s *UserService) Register(username string, password string, usertype string) (id uint, message string) {
	_, err := s.repo.FindUserByName(username)
	if err == nil {
		return 0, model.ErrUserAlreadyExists
	}
	id, err = s.repo.AddUser(username, password, usertype)
	if err != nil {
		return 0, model.ErrServerInternal
	}
	return id, ""
}

// Login 用户登录
// 用户指针为空即失败
func (s *UserService) Login(username string, password string) (user *model.User, message string) {
	targetuser, err := s.repo.FindUserByName(username)
	if err != nil {
		return nil, model.ErrUserNotFound
	}
	if targetuser.Password != password {
		return nil, model.ErrPasswordWrong
	}

	return targetuser, ""
}

// Logout 用户登出
func Logout(username string) error {
	return nil
}
