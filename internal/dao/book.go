package dao

import (
	"errors"

	"github.com/Cheerdoge/library-manage-system/internal/global"
	"github.com/Cheerdoge/library-manage-system/internal/model"
)

// FindBook 通过ID查找图书
// 成功：书指针，nil
// 失败：nil，错误信息
func FindBook(ID uint) (*model.Book, error) {
	var book model.Book
	//根据id查询，并把结果填入book
	result := global.DB.First(&book, ID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &book, nil
}

// AddBook 新增单个图书
// 成功：图书id, nil
// 失败：0, 错误信息
func AddBook(bookname string, author string, sum_num int) (uint, error) {
	var book model.Book
	book.Bookname = bookname
	book.Author = author
	book.SumNum = sum_num
	book.BorrNum = 0
	book.NowNum = sum_num
	result := global.DB.Create(&book)
	if result.Error != nil {
		return 0, result.Error
	}
	return book.ID, nil
}

// AddBooks 新增多本图书
// 成功：图书id切片, nil
// 失败：0, 错误信息
func AddBooks(newbooks []*model.Book) ([]uint, error) {
	var ids []uint
	result := global.DB.Create(&newbooks)
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
func DelBoook(book *model.Book) error {
	result := global.DB.Delete(book)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UpdateBook 更新图书信息,主要用做数量更新
// 成功：nil
// 失败：错误信息
func UpdateBook(bookid uint, sum_num int) error {
	var book *model.Book
	book, err := FindBook(bookid)
	if err != nil {
		return errors.New("书籍不存在")
	}
	book.SumNum = sum_num
	book.NowNum = sum_num - book.BorrNum
	result := global.DB.Save(&book)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetAllBooks 获取所有图书
// 成功：图书切片，nil
// 失败：nil，错误信息
func GetAllBooks() ([]model.BookInfo, error) {
	var books []model.BookInfo
	var bookModels []model.Book
	result := global.DB.Find(&bookModels)
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
func GetAvailableBooks() ([]model.BookInfo, error) {
	var books []model.BookInfo
	var bookModels []model.Book
	result := global.DB.Where("now_num > ?", 0).Find(&bookModels)
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
