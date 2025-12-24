package dao

import (
	"errors"

	"github.com/Cheerdoge/library-manage-system/internal/model"
	"gorm.io/gorm"
)

type BookDao struct {
	db *gorm.DB
}

func NewBookDao(db *gorm.DB) *BookDao {
	return &BookDao{
		db: db,
	}
}

// FindBookById 通过ID查找图书
// 成功：书指针，nil
// 失败：nil，错误信息
func (dao *BookDao) FindBookById(ID uint) (*model.BookInfo, error) {
	var book model.Book
	//根据id查询，并把结果填入book
	result := dao.db.First(&book, ID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &model.BookInfo{
		ID:       book.ID,
		Bookname: book.Bookname,
		Author:   book.Author,
		SumNum:   book.SumNum,
		BorrNum:  book.BorrNum,
		NowNum:   book.NowNum,
	}, nil
}

// FindBookByName 通过名称查找图书
// 成功：书指针，nil
// 失败：nil，错误信息
func (dao *BookDao) FindBookByName(name string) (*model.BookInfo, error) {
	var book model.Book
	//根据name查询，并把结果填入book
	result := dao.db.First(&book, "bookname = ?", name)
	if result.Error != nil {
		return nil, result.Error
	}

	return &model.BookInfo{
		ID:       book.ID,
		Bookname: book.Bookname,
		Author:   book.Author,
		BorrNum:  book.BorrNum,
		NowNum:   book.NowNum,
	}, nil
}

// AddBook 新增单个图书
// 成功：图书id, nil
// 失败：0, 错误信息
func (dao *BookDao) AddBook(bookname string, author string, sum_num int) (uint, error) {
	var book model.Book
	book.Bookname = bookname
	book.Author = author
	book.SumNum = sum_num
	book.BorrNum = 0
	book.NowNum = sum_num
	result := dao.db.Create(&book)
	if result.Error != nil {
		return 0, result.Error
	}
	return book.ID, nil
}

// AddBooks 新增多本图书
// 成功：图书id切片, nil
// 失败：0, 错误信息
func (dao *BookDao) AddBooks(newbooks []*model.Book) ([]uint, error) {
	var ids []uint
	result := dao.db.Create(&newbooks)
	if result.Error != nil {
		return nil, result.Error
	}
	for _, book := range newbooks {
		ids = append(ids, book.ID)
	}
	return ids, nil
}

// DelBook 删除图书
// 成功：nil
// 失败：错误信息
func (dao *BookDao) DelBook(bookid uint) error {
	result := dao.db.Delete(&model.Book{}, bookid)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UpdateBook 更新图书信息,主要用做数量更新
// 成功：nil
// 失败：错误信息
func (dao *BookDao) UpdateBook(bookid uint, change_num int, bor_num int, return_num int) error {
	var book *model.BookInfo
	book, err := dao.FindBookById(bookid)
	if err != nil {
		return errors.New("书籍不存在")
	}
	book.SumNum = book.SumNum + change_num
	book.BorrNum = book.BorrNum + bor_num - return_num
	book.NowNum = book.SumNum - book.BorrNum
	result := dao.db.Save(book)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetAllBooks 获取所有图书
// 成功：图书切片，nil
// 失败：nil，错误信息
func (dao *BookDao) FindAllBooks() ([]model.BookInfo, error) {
	var books []model.BookInfo
	var bookModels []model.Book
	result := dao.db.Find(&bookModels)
	if result.Error != nil {
		return nil, result.Error
	}
	for _, book := range bookModels {
		bookInfo := model.BookInfo{
			ID:       book.ID,
			Bookname: book.Bookname,
			Author:   book.Author,
			BorrNum:  book.BorrNum,
			NowNum:   book.NowNum,
		}
		books = append(books, bookInfo)
	}
	return books, nil
}

// GetAvailableBooks 获取所有可借图书
// 成功：图书切片，nil
// 失败：nil，错误信息
func (dao *BookDao) FindAvailableBooks() ([]model.BookInfo, error) {
	var books []model.BookInfo
	var bookModels []model.Book
	result := dao.db.Where("now_num > ?", 0).Find(&bookModels)
	if result.Error != nil {
		return nil, result.Error
	}
	for _, book := range bookModels {
		bookInfo := model.BookInfo{
			ID:       book.ID,
			Bookname: book.Bookname,
			Author:   book.Author,
			BorrNum:  book.BorrNum,
			NowNum:   book.NowNum,
		}
		books = append(books, bookInfo)

	}
	return books, nil
}

// ModifyStore 配合事务的库存修改
func (dao *BookDao) ModifyStore(tx *gorm.DB, bookid uint, num int) error {
	result := tx.Model(&model.Book{}).Where("id = ?", bookid).Update("now_num", gorm.Expr("now_num + ?", num))
	if result.Error != nil {
		return result.Error
	}
	return nil
}
