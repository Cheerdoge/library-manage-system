package handler

import (
	"strconv"

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

// GetBooksHandler 获取所有书籍信息，支持查询可借图书
func (h *BookHandler) GetBooksHandler(c *gin.Context) {
	available := c.Query("available")
	name := c.Query("name")

	if available != "true" && available != "" {
		web.FailWithMessage(c, "available参数有误")
		return
	}

	if available == "true" && name != "" {
		book, msg := h.bookservice.GetAvailableBooksByName(name)
		if msg != "" {
			web.FailWithMessage(c, msg)
			return
		}
		web.OkWithData(c, []interface{}{book})
		return
	}

	if available == "true" {
		booklist, msg := h.bookservice.GetAvailableBooks()
		if msg != "" {
			web.FailWithMessage(c, msg)
			return
		}
		web.OkWithData(c, booklist)
		return
	}

	if name != "" {
		book, msg := h.bookservice.GetBookByName(name)
		if msg != "" {
			web.FailWithMessage(c, msg)
			return
		}
		// 为了保持返回格式统一为数组，建议包一层
		web.OkWithData(c, []interface{}{book})
		return
	}

	booklist, msg := h.bookservice.GetAllBooks()
	if msg != "" {
		web.FailWithMessage(c, msg)
		return
	}
	web.OkWithData(c, booklist)
}

// GetBookByIdHandler 通过ID获取书籍信息
func (h *BookHandler) GetBookByIdHandler(c *gin.Context) {
	// var req web.FindBookById
	// err := c.ShouldBindJSON(&req)
	// if err != nil {
	// 	web.FailWithMessage(c, "请求参数有误")
	// 	return
	// }
	stringid := c.Param("id")
	id, err := strconv.ParseUint(stringid, 10, 0)
	if err != nil {
		web.FailWithMessage(c, "请求参数有误")
		return
	}
	book, msg := h.bookservice.GetBookById(uint(id))
	if msg != "" {
		web.FailWithMessage(c, msg)
		return
	}
	web.OkWithData(c, book)
}

// DeleteBookHandler 删除图书
func (h *BookHandler) DeleteBookHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		web.FailWithMessage(c, "ID参数错误")
		return
	}

	ok, msg := h.bookservice.RemoveBook(uint(id))
	if !ok {
		web.FailWithMessage(c, msg)
		return
	}
	web.OkWithMessage(c, "删除图书成功")
}

// 新增图书 AddBookHandler 返回新增图书的id
func (h *BookHandler) AddBookHandler(c *gin.Context) {
	var req web.AddBook
	err := c.ShouldBindJSON(&req)
	if err != nil {
		web.FailWithMessage(c, "请求参数有误")
		return
	}
	if req.Bookname == "" || req.Author == "" || req.Num <= 0 {
		web.FailWithMessage(c, "图书名称、作者不能为空,数量需大于0")
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

// UpdateBookHandler 更新图书数量
func (h *BookHandler) UpdateBookHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		web.FailWithMessage(c, "ID参数错误")
		return
	}

	var req web.UpdateBook
	if err := c.ShouldBindJSON(&req); err != nil {
		web.FailWithMessage(c, "请求参数有误")
		return
	}

	ok, msg := h.bookservice.ModifyBook(uint(id), req.Num, 0, 0)
	if !ok {
		web.FailWithMessage(c, msg)
		return
	}
	web.OkWithMessage(c, "更新图书数量成功")
}
