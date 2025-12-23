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

// Adduser添加新用户
// 成功：用户id，nil
// 失败：0，错误信息
func (dao *UserDao) Adduser(username string, password string, isadmin bool) (uint, error) {
	var user model.User
	user.UserName = username
	user.Password = password
	user.IsAdmin = isadmin
	result := dao.db.Create(&user)
	if result.Error != nil {
		return 0, result.Error
	}
	return user.ID, nil
}

// FindUserById 通过ID查找用户
// 成功：用户指针，nil
// 失败：nil，错误信息
func (dao *UserDao) FindUserById(Id uint) (*model.User, error) {
	var user model.User
	result := dao.db.First(&user, Id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// FindUserByName 通过名称查找用户
// 成功：用户指针，nil
// 失败：nil，错误信息
func (dao *UserDao) FindUserByName(name string) (*model.User, error) {
	var user model.User
	result := dao.db.Where("name = ?", name).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// UpdatePassword 更新用户密码
// 管理员直接操作，用户验证密码
// 成功：nil
// 失败：错误信息
func (dao *UserDao) UpdatePassword(userid uint, newpassword string) error {
	var user *model.User
	user, err := dao.FindUserById(userid)
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

// UpdataUserInfo 更新用户信息
// 成功：nil
// 失败：错误信息
func (dao *UserDao) UpdateUserInfo(userid uint, username string, telenum string) error {
	var user *model.User
	user, err := dao.FindUserById(userid)
	if err != nil {
		return errors.New("用户不存在")
	}
	user.UserName = username
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
	user, err := dao.FindUserById(userid)
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
			UserName:   user.UserName,
			Telenum:    user.Telenum,
			IsAdmin:    user.IsAdmin,
			BorrRecNum: user.NowBorrNum,
		}
		userinfos = append(userinfos, userinfo)
	}
	return userinfos, nil
}
