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
// 成功：应还日期，nil
// 失败：空，错误信息
func (dao *BorrowRecordDao) BorrowBook(bookid uint, userid uint) (time.Time, error) {
	var record model.BorrowRecord
	record.BookID = bookid
	record.UserID = userid
	record.BorrowDate = time.Now()
	record.ReturnDate = time.Time{}
	record.ShouldReturn = record.BorrowDate.AddDate(0, 0, 7) // 测试借书期限为7天
	record.State = "borrowing"
	result := dao.db.Create(&record)
	if result.Error != nil {
		return time.Time{}, result.Error
	}
	return record.ShouldReturn, nil
}

// ReturnBook 还书
// 成功：nil, 真为守时，假为逾时
// 失败：错误信息, 假
func (dao *BorrowRecordDao) ReturnBook(recordid uint) (error, bool) {
	var record model.BorrowRecord
	result := dao.db.First(&record, recordid)
	if result.Error != nil {
		return result.Error, false
	}
	record.ReturnDate = time.Now()
	record.State = "returned"
	result = dao.db.Save(&record)
	if result.Error != nil {
		return result.Error, false
	}
	if record.ReturnDate.After(record.ShouldReturn) {
		return nil, false
	}
	return nil, true
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
