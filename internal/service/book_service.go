package service

import (
	"github.com/Cheerdoge/library-manage-system/internal/model"
)

type BookRepository interface {
	FindBookById(ID uint) (*model.Book, error)
	FindBookByName(name string) (*model.Book, error)
	AddBook(bookname string, author string, sum_num int) (uint, error)
	AddBooks(newbooks []*model.Book) ([]uint, error)
	DelBook(book *model.Book) error
	UpdateBook(bookid uint, sum_num int) error
	GetAllBooks() ([]model.Book, error)
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
func (s *BookService) GetAllBooks() (booklist []model.Book, message string) {
	booklist, err := s.repo.GetAllBooks()
	if err != nil {
		return nil, "获取图书列表失败:" + err.Error()
	}
	return booklist, ""
}

func (s *BookService) FindBookById(ID uint) (book *model.Book, message string) {
	book, err := s.repo.FindBookById(ID)
	if err != nil {
		return nil, "查找图书失败:" + err.Error()
	}
	return book, ""
}

func (s *BookService) FindBookByName(name string) (book *model.Book, message string) {
	book, err := s.repo.FindBookByName(name)
	if err != nil {
		return nil, "查找图书失败:" + err.Error()
	}
	return book, ""
}

func (s *BookService) FindAvailableBookByName(name string) (book *model.Book, message string) {
	book, err := s.repo.FindBookByName(name)
	if err != nil {
		return nil, "查找可借阅图书失败:" + err.Error()
	}
	if book.NowNum <= 0 {
		return nil, "图书暂无库存"
	}
	return book, ""
}
