package dao

import (
	"fmt"
	"time"

	"github.com/Cheerdoge/library-manage-system/internal/model"
	"gorm.io/gorm"
)

type SessionDao struct {
	db *gorm.DB
}

func NewSessionDao(db *gorm.DB) *SessionDao {
	return &SessionDao{
		db: db,
	}
}

// AddSession 创建新session
// 返回token,有效期30分钟
func (dao *SessionDao) AddSession(userID uint, username string, IsAdmin bool) (token string, err error) {
	var session model.Session
	session.UserID = userID
	session.UserName = username
	session.IsAdmin = IsAdmin
	token = fmt.Sprintf("session-%d-%d", userID, time.Now().UnixNano())
	session.ExpiresAt = time.Now().Add(30 * time.Minute)
	session.Token = token
	result := dao.db.Create(&session)
	if result.Error != nil {
		return "", result.Error
	}
	return token, nil
}

// FindSessionByToken 通过token获取session
func (dao *SessionDao) FindSessionByToken(token string) (*model.Session, error) {
	var session model.Session
	result := dao.db.Where("token = ?", token).First(&session)
	if result.Error != nil {
		return nil, result.Error
	}
	return &session, nil
}

// DeleteSessionByToken 通过token删除session
func (dao *SessionDao) DeleteSessionByToken(token string) error {
	result := dao.db.Where("token = ?", token).Delete(&model.Session{})
	return result.Error
}

func (dao *SessionDao) FindSessionByUserId(userId uint) (*model.Session, error) {
	var session model.Session
	result := dao.db.First(&session, "user_id = ?", userId)
	if result.Error != nil {
		return nil, result.Error
	}
	return &session, nil
}
