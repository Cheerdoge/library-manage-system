package service

import (
	"time"

	"github.com/Cheerdoge/library-manage-system/internal/model"
	"gorm.io/gorm"
)

type BorrowRepository interface {
	Create(db *gorm.DB, record *model.BorrowRecord) error
	ReturnBook(db *gorm.DB, recordid uint) error
	FindBorrowRecord(userid uint) ([]model.BorrowRecord, error)
	FindAllBorrowRecords() ([]model.BorrowRecord, error)
	FindBorrowRecordById(recordid uint) (*model.BorrowRecord, error)
}

type BorrowService struct {
	db          *gorm.DB
	borrowrepo  BorrowRepository
	userservice *UserService
	bookservice *BookService
}

func NewBorrowService(db *gorm.DB, repo BorrowRepository, userservice *UserService, bookservice *BookService) *BorrowService {
	return &BorrowService{
		db:          db,
		borrowrepo:  repo,
		userservice: userservice,
		bookservice: bookservice,
	}
}

// Borrow 借书
// 失败返回空和错误信息
func (s *BorrowService) Borrow(bookid uint, userid uint, booknum int) (shouldreturn time.Time, message string) {
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
	book, err := s.bookservice.repo.FindBookById(bookid)
	if err != nil {
		return time.Time{}, "图书不存在:" + err.Error()
	}
	if book.NowNum < booknum {
		return time.Time{}, "图书库存不足"
	}
	shouldreturn = time.Now().AddDate(0, 0, 7) // 测试借书期限为7天
	record := &model.BorrowRecord{
		BookID:       bookid,
		UserID:       userid,
		BookNum:      booknum,
		BorrowDate:   time.Now(),
		ShouldReturn: shouldreturn,
		State:        "borrowing",
	}
	err = s.db.Transaction(func(tx *gorm.DB) error {
		err := s.bookservice.repo.ModifyStore(tx, bookid, -booknum)
		if err != nil {
			return err
		}
		err = s.borrowrepo.Create(tx, record)
		if err != nil {
			return err
		}
		err = s.userservice.userrepo.ModifyUserNum(tx, userid, booknum, 0)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return time.Time{}, "借书失败:" + err.Error()
	}
	return shouldreturn, ""
}

// Return 还书
// 成功：真为守时，假为逾时，“”,判定是否有错误信息来判断是否操作成功
func (s *BorrowService) Return(recordid uint) (isontime bool, message string) {
	var targetrecord *model.BorrowRecord
	targetrecord, err := s.borrowrepo.FindBorrowRecordById(recordid)
	if err != nil {
		return false, "查找借书记录失败:" + err.Error()
	}
	if targetrecord == nil {
		return false, "未找到对应的借书记录"
	}
	err = s.db.Transaction(func(tx *gorm.DB) error {
		err := s.borrowrepo.ReturnBook(tx, targetrecord.ID)
		if err != nil {
			return err
		}
		err = s.bookservice.repo.ModifyStore(tx, targetrecord.BookID, targetrecord.BookNum)
		if err != nil {
			return err
		}
		if time.Now().After(targetrecord.ShouldReturn) {
			err = s.userservice.userrepo.ModifyUserNum(tx, targetrecord.UserID, -targetrecord.BookNum, 1)
			isontime = false
		} else {
			err = s.userservice.userrepo.ModifyUserNum(tx, targetrecord.UserID, -targetrecord.BookNum, 0)
			isontime = true
		}
		if err != nil {
			return err
		}
		return nil
	})
	return isontime, ""
}

// GetUserBorrowRecords 获取用户未归还的借书记录
// 失败返回nil，错误信息
func (s *BorrowService) GetUserBorrowRecords(userid uint) (records []model.BorrowRecordInfo, message string) {
	_, err := s.userservice.userrepo.FindUserById(userid)
	if err != nil {
		return nil, model.ErrUserNotFound
	}
	recordlist, err := s.borrowrepo.FindBorrowRecord(userid)
	for _, record := range recordlist {
		recordInfo := model.BorrowRecordInfo{
			ID:           record.ID,
			BookID:       record.BookID,
			UserID:       record.UserID,
			BookNum:      record.BookNum,
			BorrowDate:   record.BorrowDate.Format("2006-01-02 15:04:05"),
			ShouldReturn: record.ShouldReturn.Format("2006-01-02 15:04:05"),
			State:        record.State,
		}
		records = append(records, recordInfo)
	}
	if err != nil {
		return nil, err.Error()
	}
	return records, ""
}

// GetAllBorrowRecords 获取所有状态为未归还的借书记录,管理员专属
// 成功：借书记录切片，“”
// 失败：nil，错误信息
func (s *BorrowService) GetAllBorrowRecords(isadmin bool) (records []model.BorrowRecordInfo, message string) {
	if !isadmin {
		return nil, model.ErrForbidden
	}
	recordlist, err := s.borrowrepo.FindAllBorrowRecords()
	if err != nil {
		return nil, err.Error()
	}
	for _, record := range recordlist {
		recordInfo := model.BorrowRecordInfo{
			ID:           record.ID,
			BookID:       record.BookID,
			UserID:       record.UserID,
			BookNum:      record.BookNum,
			BorrowDate:   record.BorrowDate.Format("2006-01-02 15:04:05"),
			ShouldReturn: record.ShouldReturn.Format("2006-01-02 15:04:05"),
			State:        record.State,
		}
		records = append(records, recordInfo)
	}
	return records, ""
}
