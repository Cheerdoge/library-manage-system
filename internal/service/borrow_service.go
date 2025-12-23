package service

import (
	"time"

	"github.com/Cheerdoge/library-manage-system/internal/model"
)

type BorrowRepository interface {
	BorrowBook(bookid uint, userid uint) (time.Time, error)
	ReturnBook(recordid uint) (error, bool)
	FindBorrowRecord(userid uint) ([]model.BorrowRecord, error)
	FindAllBorrowRecords() ([]model.BorrowRecord, error)
}

type BorrowService struct {
	borrowrepo  BorrowRepository
	userservice *UserService
	bookservice *BookService
}

func NewBorrowService(repo BorrowRepository, userservice *UserService, bookservice *BookService) *BorrowService {
	return &BorrowService{
		borrowrepo:  repo,
		userservice: userservice,
		bookservice: bookservice,
	}
}

// Borrow 借书
// 失败返回空和错误信息
func (s *BorrowService) Borrow(bookid uint, userid uint) (shouldreturn time.Time, message string) {
	user, err := s.userservice.userrepo.FindUserById(userid)
	if err != nil {
		return time.Time{}, model.ErrUserNotFound
	}
	if user.NowBorrNum >= 5 {
		return time.Time{}, model.ErrBorrowLimitExceeded
	}
	if user.OverdueNum > 3 {
		return time.Time{}, model.ErrUserOverdueLimitExceeded
	}
	shouldreturn, err = s.borrowrepo.BorrowBook(bookid, userid)
	if err != nil {
		return time.Time{}, err.Error()
	}
	return shouldreturn, ""
}

// Return 还书
// 成功：真为守时，假为逾时，“”
func (s *BorrowService) Return(userid uint, bookid uint) (isontime bool, message string) {
	user, err := s.userservice.userrepo.FindUserById(userid)
	if err != nil {
		return false, model.ErrUserNotFound
	}
	records, err := s.borrowrepo.FindBorrowRecord(userid)
	if err != nil {
		return false, "获取借书记录失败:" + err.Error()
	}
	var targetrecord *model.BorrowRecord
	for _, record := range records {
		if record.BookID == bookid && record.State == "borrowing" {
			targetrecord = &record
			break
		}
	}
	if targetrecord == nil {
		return false, "未找到对应的借书记录"
	}
	err, isontime = s.borrowrepo.ReturnBook(targetrecord.ID)
	if err != nil {
		return false, err.Error()
	}
	if !isontime {
		user.OverdueNum++
		s.userservice.userrepo.UpdateUserInfo(user.ID, user.UserName, user.Telenum, user.OverdueNum)
	}
	return isontime, ""
}

// GetUserBorrowRecords 获取用户未归还的借书记录
// 失败返回nil，错误信息
func (s *BorrowService) GetUserBorrowRecords(userid uint) (records []model.BorrowRecord, message string) {
	_, err := s.userservice.userrepo.FindUserById(userid)
	if err != nil {
		return nil, model.ErrUserNotFound
	}
	records, err = s.borrowrepo.FindBorrowRecord(userid)
	if err != nil {
		return nil, err.Error()
	}
	return records, ""
}

// GetAllBorrowRecords 获取所有用户未归还的借书记录,管理员专属
// 成功：借书记录切片，“”
// 失败：nil，错误信息
func (s *BorrowService) GetAllBorrowRecords(isadmin bool) (records []model.BorrowRecord, message string) {
	if !isadmin {
		return nil, model.ErrForbidden
	}
	records, err := s.borrowrepo.FindAllBorrowRecords()
	if err != nil {
		return nil, err.Error()
	}
	return records, ""
}
