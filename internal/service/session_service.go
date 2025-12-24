package service

import (
	"time"

	"github.com/Cheerdoge/library-manage-system/internal/model"
)

type SessionRepository interface {
	AddSession(userID uint, username string, IsAdmin bool) (token string, err error)
	FindSessionByToken(token string) (*model.Session, error)
	DeleteSessionByToken(token string) error
	FindSessionByUserId(userId uint) (*model.Session, error)
}

type SessionService struct {
	repo SessionRepository
}

func NewSessionService(repo SessionRepository) *SessionService {
	return &SessionService{
		repo: repo,
	}
}

// CreateSession 创建新session并返回token
func (s *SessionService) CreateSession(userId uint, username string, IsAdmin bool) (token string, messsage string) {
	token, err := s.repo.AddSession(userId, username, IsAdmin)
	if err != nil {
		return "", "创建session失败:" + err.Error()
	}
	return token, ""
}

// CheckeSessionByToken 通过token获取session并验证是否过期
func (s *SessionService) CheckSessionByToken(token string) (session *model.Session, message string) {
	session, err := s.repo.FindSessionByToken(token)
	if err != nil {
		return nil, "获取session失败:" + err.Error()
	}
	if session.ExpiresAt.Before(time.Now()) {
		s.repo.DeleteSessionByToken(token)
		return nil, "session已过期"
	}
	return session, ""
}

// DelSessionByToken 通过token删除session
func (s *SessionService) DelSessionByToken(token string) (message string) {
	_, err := s.repo.FindSessionByToken(token)
	if err != nil {
		return "session不存在:" + err.Error()
	}
	err = s.repo.DeleteSessionByToken(token)
	if err != nil {
		return "删除session失败:" + err.Error()
	}
	return ""
}
