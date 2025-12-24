package dao

import (
	"errors"
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
func (dao *BorrowRecordDao) Create(db *gorm.DB, record *model.BorrowRecord) error {
	return db.Create(record).Error
}

// ReturnBook 还书
// 成功：真，nil
func (dao *BorrowRecordDao) ReturnBook(db *gorm.DB, recordid uint) error {
	result := db.Model(&model.BorrowRecord{}).
		Where("id = ? AND state = ?", recordid, "borrowing").
		Updates(map[string]interface{}{
			"state":     "returned",
			"return_at": time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("未找到对应的借书记录或该记录已归还")
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

// FindBorrowRecordById 通过借书记录ID查找借书记录
func (dao *BorrowRecordDao) FindBorrowRecordById(recordid uint) (*model.BorrowRecord, error) {
	var record model.BorrowRecord
	result := dao.db.First(&record, "id = ?", recordid)
	if result.Error != nil {
		return nil, result.Error
	}
	return &record, nil
}
