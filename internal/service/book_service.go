package service

import (
	"github.com/Cheerdoge/library-manage-system/internal/model"
	"gorm.io/gorm"
)

type BookRepository interface {
	FindBookById(ID uint) (*model.BookInfo, error)
	FindBookByName(name string) ([]model.BookInfo, error)
	AddBook(bookname string, author string, sum_num int) (uint, error)
	AddBooks(newbooks []*model.Book) ([]uint, error)
	DelBook(bookid uint) error
	UpdateBook(bookid uint, change_num int, bor_num int, return_num int) error
	FindAllBooks() ([]model.BookInfo, error)
	FindAvailableBooks() ([]model.BookInfo, error)
	ModifyStore(tx *gorm.DB, bookid uint, nownum int, bornum int) error
}

type BookService struct {
	repo BookRepository
}

func NewBookService(repo BookRepository) *BookService {
	return &BookService{
		repo: repo,
	}
}

// GetAllBooks 获取所有图书
// 图书列表为空则失败
func (s *BookService) GetAllBooks() (booklist []model.BookInfo, message string) {
	booklist, err := s.repo.FindAllBooks()
	if err != nil {
		return nil, "获取图书列表失败:" + err.Error()
	}
	return booklist, ""
}

// GetBookById 通过Id查找图书
func (s *BookService) GetBookById(ID uint) (book *model.BookInfo, message string) {
	book, err := s.repo.FindBookById(ID)
	if err != nil {
		return nil, "图书不存在:" + err.Error()
	}
	return book, ""
}

// GetBookByName 通过名称查找图书
func (s *BookService) GetBookByName(name string) (books []model.BookInfo, message string) {
	books, err := s.repo.FindBookByName(name)
	if err != nil {
		return nil, "查询失败:" + err.Error()
	}
	if len(books) == 0 {
		return nil, "图书不存在"
	}
	return books, ""
}

// GetAvailableBooksByName 通过名称查找可借阅图书
func (s *BookService) GetAvailableBooksByName(name string) (books []model.BookInfo, message string) {
	books, err := s.repo.FindBookByName(name)
	if err != nil {
		return nil, "查找可借阅图书失败:" + err.Error()
	}
	if len(books) == 0 {
		return nil, "图书不存在"
	}
	var availableBooks []model.BookInfo
	for _, book := range books {
		if book.NowNum > 0 {
			availableBooks = append(availableBooks, book)
		}
	}
	if len(availableBooks) == 0 {
		return nil, "所查询图书暂无可借阅库存"
	}
	return availableBooks, ""
}

// GetAvailableBooks 获取所有可借图书
func (s *BookService) GetAvailableBooks() (booklist []model.BookInfo, message string) {
	booklist, err := s.repo.FindAvailableBooks()
	if err != nil {
		return nil, "获取可借图书列表失败:" + err.Error()
	}
	return booklist, ""
}

// NewBook 新增图书
func (s *BookService) NewBook(bookname string, author string, sum_num int) (bookid uint, message string) {
	bookid, err := s.repo.AddBook(bookname, author, sum_num)
	if err != nil {
		return 0, "新增图书失败:" + err.Error()
	}
	return bookid, ""
}

// // NewBooks 新增多本图书
// func (s *BookService) NewBooks(newbooks []*model.Book) (bookidlist []uint, message string) {
// 	bookidlist, err := s.repo.AddBooks(newbooks)
// 	if err != nil {
// 		return nil, "新增图书失败:" + err.Error()
// 	}

// 	return bookidlist, ""
// }

// RemoveBook 删除图书
// bool值表示是否成功
func (s *BookService) RemoveBook(bookid uint) (ok bool, message string) {
	err := s.repo.DelBook(bookid)
	if err != nil {
		return false, "删除图书失败:" + err.Error()
	}
	return true, ""
}

// ModifyBook 更新图书信息,主要用做数量更新
// bool值表示是否成功
func (s *BookService) ModifyBook(bookid uint, change_num int, bor_num int, return_num int) (ok bool, message string) {
	err := s.repo.UpdateBook(bookid, change_num, bor_num, return_num)
	if err != nil {
		return false, "更新图书信息失败:" + err.Error()
	}
	return true, ""
}
