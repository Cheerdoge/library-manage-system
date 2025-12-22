package dao

import (
	"errors"

	"github.com/Cheerdoge/library-manage-system/internal/model"
	"gorm.io/gorm"
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}

// Register注册新用户
// 成功：用户id，nil
// 失败：0，错误信息
func (dao *UserDao) Register(username string, password string, usertype string) (uint, error) {
	var user model.User
	user.Name = username
	user.Password = password
	user.Type = usertype
	result := dao.db.Create(&user)
	if result.Error != nil {
		return 0, result.Error
	}
	return user.ID, nil
}

// FindUser 通过ID查找用户
// 成功：用户指针，nil
// 失败：nil，错误信息
func (dao *UserDao) FindUser(Id uint) (*model.User, error) {
	var user model.User
	result := dao.db.First(&user, Id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// Login 用户登录
// 成功：用户指针，nil
// 失败：nil，错误信息
func (dao *UserDao) Login(username string, password string) (*model.User, error) {
	var user model.User
	result := dao.db.Where("name = ? AND password = ?", username, password).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// Logout 用户登出
// 成功：nil
// 失败：错误信息
func Logout(username string) error {
	return nil
}

// ChangePassword 修改用户密码
// 管理员直接操作，用户验证密码
// 成功：nil
// 失败：错误信息
func (dao *UserDao) ChangePassword(userid uint, newpassword string) error {
	var user *model.User
	user, err := dao.FindUser(userid)
	if err != nil {
		return err
	}
	user.Password = newpassword
	result := dao.db.Save(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// ChangeUserInfo 修改用户信息
// 成功：nil
// 失败：错误信息
func (dao *UserDao) ChangeUserInfo(userid uint, username string, telenum string) error {
	var user *model.User
	user, err := dao.FindUser(userid)
	if err != nil {
		return errors.New("用户不存在")
	}
	user.Name = username
	user.Telenum = telenum
	result := dao.db.Save(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// DeleUser 删除用户
// 管理员直接操作，用户验证密码
// 成功：nil
// 失败：错误信息
func (dao *UserDao) DeleUser(userid uint) error {
	var user *model.User
	user, err := dao.FindUser(userid)
	if err != nil {
		return errors.New("用户不存在")
	}
	result := dao.db.Delete(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetAllUsers 获取所有用户信息
// 仅管理员可用
// 成功：用户信息切片, nil
// 失败：nil, 错误信息
func (dao *UserDao) GetAllUsers() ([]model.UserInfo, error) {
	var users []model.User
	var userinfos []model.UserInfo
	result := dao.db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, user := range users {
		userinfo := model.UserInfo{
			ID:         user.ID,
			Name:       user.Name,
			Telenum:    user.Telenum,
			Type:       user.Type,
			BorrRecNum: user.NowBorrNum,
		}
		userinfos = append(userinfos, userinfo)
	}
	return userinfos, nil
}
