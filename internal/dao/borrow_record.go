package dao

import (
	"time"

	"github.com/Cheerdoge/library-manage-system/internal/model"
	"gorm.io/gorm"
)

type BorrowRecordDao struct {
	db *gorm.DB
}

func NewBorrowRecordDao(db *gorm.DB) *BorrowRecordDao {
	return &BorrowRecordDao{
		db: db,
	}
}

// BorrowBook 借书
// 成功：借书记录id，nil
// 失败：0，错误信息
func (dao *BorrowRecordDao) BorrowBook(bookid uint, userid uint, borrowedAt time.Time) (time.Time, error) {
	var record model.BorrowRecord
	record.BookID = bookid
	record.UserID = userid
	record.BorrowDate = borrowedAt
	record.State = "borrowing"
	result := dao.db.Create(&record)
	if result.Error != nil {
		return time.Time{}, result.Error
	}
	return record.BorrowDate, nil
}

// ReturnBook 还书
// 成功：nil
// 失败：错误信息
func (dao *BorrowRecordDao) ReturnBook(recordid uint, returnedAt time.Time) error {
	var record model.BorrowRecord
	result := dao.db.First(&record, recordid)
	if result.Error != nil {
		return result.Error
	}
	record.ReturnDate = returnedAt
	record.State = "returned"
	result = dao.db.Save(&record)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// FindBorrowRecord 通过ID查找某用户未归还的借书记录
// 成功：借书记录切片，nil
// 失败：nil，错误信息
func (dao *BorrowRecordDao) FindBorrowRecord(userid uint) ([]model.BorrowRecord, error) {
	var records []model.BorrowRecord
	result := dao.db.Find(&records, "user_id = ? AND state = ?", userid, "borrowing")
	if result.Error != nil {
		return nil, result.Error
	}
	return records, nil
}

// FindAllBorrowRecords 查找所有未归还借书记录
// 成功：借书记录切片，nil
// 失败：nil，错误信息
func (dao *BorrowRecordDao) FindAllBorrowRecords() ([]model.BorrowRecord, error) {
	var records []model.BorrowRecord
	result := dao.db.Find(&records, "state = ?", "borrowing")
	if result.Error != nil {
		return nil, result.Error
	}
	return records, nil
}
