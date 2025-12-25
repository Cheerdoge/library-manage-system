package handler

import (
	"strconv"

	"github.com/Cheerdoge/library-manage-system/internal/model"
	"github.com/Cheerdoge/library-manage-system/internal/service"
	"github.com/Cheerdoge/library-manage-system/web"
	"github.com/gin-gonic/gin"
)

type BorrowHandler struct {
	borrowservice *service.BorrowService
}

func NewBorrowHandler(borrowservice *service.BorrowService) *BorrowHandler {
	return &BorrowHandler{
		borrowservice: borrowservice,
	}
}

// 请求借书 BorrowHandler
func (h *BorrowHandler) BorrowBookHandler(c *gin.Context) {
	var req web.BorrowBook
	err := c.ShouldBindJSON(&req)
	if err != nil {
		web.FailWithMessage(c, "请求参数有误")
	}
	principal, err := model.GetPrincipal(c)
	if err != nil {
		web.FailWithMessage(c, "无法获取用户信息")
		return
	}
	shouldreturn, msg := h.borrowservice.Borrow(req.Bookid, principal.UserID, req.BookNum)
	if msg != "" {
		web.FailWithMessage(c, msg)
		return
	}
	web.Ok(c, "借书成功,应还时间:", gin.H{"should_return": shouldreturn.Format("2006-01-02")})
}

// 请求还书 ReturnBookHandler
func (h *BorrowHandler) ReturnBookHandler(c *gin.Context) {
	recordidstr := c.Param("recordid")
	recordid, err := strconv.Atoi(recordidstr)
	if err != nil {
		web.FailWithMessage(c, "请求参数有误")
		return
	}
	isontime, msg := h.borrowservice.Return(uint(recordid))
	if msg != "" {
		web.FailWithMessage(c, msg)
		return
	}
	if isontime {
		web.OkWithMessage(c, "还书成功，守时归还")
	} else {
		web.OkWithMessage(c, "还书成功，逾时归还，请注意下次守时")
	}
}

// 获取用户未归还借书记录 GetUserBorrowRecordsHandler
func (h *BorrowHandler) GetUserBorrowRecordsHandler(c *gin.Context) {
	principal, err := model.GetPrincipal(c)
	if err != nil {
		web.FailWithMessage(c, "无法获取用户信息")
		return
	}
	records, msg := h.borrowservice.GetUserBorrowRecords(principal.UserID)
	if msg != "" {
		web.FailWithMessage(c, msg)
		return
	}
	if len(records) == 0 {
		web.OkWithMessage(c, "暂无未归还的借书记录")
		return
	}
	web.OkWithData(c, records)
}

// 管理员专属，获取所有未归还借书记录 GetAllBorrowRecordsHandler
func (h *BorrowHandler) GetAllBorrowRecordsHandler(c *gin.Context) {
	principal, err := model.GetPrincipal(c)
	if err != nil {
		web.FailWithMessage(c, "无法获取用户信息")
		return
	}
	records, msg := h.borrowservice.GetAllBorrowRecords(principal.IsAdmin)
	if msg != "" {
		web.FailWithMessage(c, msg)
		return
	}
	web.OkWithData(c, records)
}
