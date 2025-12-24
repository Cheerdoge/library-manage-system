package handler

import (
	"github.com/Cheerdoge/library-manage-system/internal/service"
	"github.com/Cheerdoge/library-manage-system/web"
	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	bookservice *service.BookService
}

func NewBookHandler(bookservice *service.BookService) *BookHandler {
	return &BookHandler{
		bookservice: bookservice,
	}
}

// GetBooksHandler 获取所有书籍信息
func (h *BookHandler) GetBooksHandler(c *gin.Context) {
	booklist, msg := h.bookservice.GetAllBooks()
	if msg != "" {
		web.FailWithMessage(c, msg)
		return
	}
	web.OkWithData(c, booklist)
}

// GetBookByNameHandler 通过name获取书籍信息
func (h *BookHandler) GetBookByNameHandler(c *gin.Context) {
	var req web.FindBook
	err := c.ShouldBindJSON(&req)
	if err != nil {
		web.FailWithMessage(c, "请求参数有误")
		return
	}
	book, msg := h.bookservice.GetBookByName(req.Bookname)
	if msg != "" {
		web.FailWithMessage(c, msg)
		return
	}
	web.OkWithData(c, book)
}

// GetBookByIdHandler 通过ID获取书籍信息 真的会有人用ID查书吗
func (h *BookHandler) GetBookByIdHandler(c *gin.Context) {
	var req web.FindBookById
	err := c.ShouldBindJSON(&req)
	if err != nil {
		web.FailWithMessage(c, "请求参数有误")
		return
	}
	book, msg := h.bookservice.GetBookById(req.BookId)
	if msg != "" {
		web.FailWithMessage(c, msg)
		return
	}
	web.OkWithData(c, book)
}

// DeleteBookHandler 删除图书
func (h *BookHandler) DeleteBookHandler(c *gin.Context) {
	var req web.DelBook
	err := c.ShouldBindJSON(&req)
	if err != nil {
		web.FailWithMessage(c, "请求参数有误")
		return
	}
	ok, msg := h.bookservice.RemoveBook(req.Id)
	if !ok {
		web.FailWithMessage(c, msg)
		return
	}
	web.OkWithMessage(c, "删除图书成功")
}

// GetAvailableBooksHandler 获取所有可借图书信息
func (h *BookHandler) GetAvailableBooksHandler(c *gin.Context) {
	booklist, msg := h.bookservice.GetAvailableBooks()
	if msg != "" {
		web.FailWithMessage(c, msg)
		return
	}
	web.OkWithData(c, booklist)
}

// 新增图书 AddBookHandler 返回新增图书的id
func (h *BookHandler) AddBookHandler(c *gin.Context) {
	var req web.AddBook
	err := c.ShouldBindJSON(&req)
	if err != nil {
		web.FailWithMessage(c, "请求参数有误")
		return
	}
	bookid, msg := h.bookservice.NewBook(req.Bookname, req.Author, req.Num)
	if msg != "" {
		web.FailWithMessage(c, msg)
		return
	}
	web.OkWithData(c, gin.H{"bookid": bookid})
}

// // 新增很多图书 AddBooksHandler 返回新增图书的id列表
// func (h *BookHandler) AddBooksHandler(c *gin.Context) {
// 	var reqs []web.AddBook
// 	err := c.ShouldBindJSON(&reqs)
// 	if err != nil {
// 		web.FailWithMessage(c, "请求参数有误")
// 		return
// 	}
// 	var newbooks []*model.Book
// 	for _, req := range reqs {
// 		newbook := &model.Book{
// 			Bookname: req.Bookname,
// 			Author:   req.Author,
// 			Num:      req.Num,
// 		}
// 		newbooks = append(newbooks, newbook)
// 	}
// 	bookids, msg := h.bookservice.NewBooks(newbooks)
// 	if msg != "" {
// 		web.FailWithMessage(c, msg)
// 		return
// 	}
// 	web.OkWithData(c, gin.H{"bookids": bookids})
// }

// 更新图书数量 UpdateBookHandler
func (h *BookHandler) UpdateBookHandler(c *gin.Context) {
	var req web.UpdateBook
	err := c.ShouldBindJSON(&req)
	if err != nil {
		web.FailWithMessage(c, "请求参数有误")
		return
	}
	ok, msg := h.bookservice.ModifyBook(req.Id, req.Num, 0, 0)
	if !ok {
		web.FailWithMessage(c, msg)
		return
	}
	web.OkWithMessage(c, "更新图书数量成功")
}
