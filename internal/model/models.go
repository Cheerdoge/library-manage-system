package model

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	ID        uint   `json:"id"`
	Bookname  string `json:"bookname"`
	Author    string `json:"author"`
	SumNum    int    `json:"sum_num"`
	BorrNum   int    `json:"borr_num"`
	NowNum    int    `json:"now_num"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type User struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Password   string `json:"-"`
	Telenum    string `json:"telenum"`
	NowBorrNum int    `json:"now_borr_num"`
	OverdueNum int    `json:"overdue_num"`
	Type       string `json:"type"` //user为普通用户，admin为管理员
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

type BorrowRecord struct {
	ID         uint `json:"id"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	BookID     uint           `json:"book_id" gorm:"column:book_id"`
	UserID     uint           `json:"user_id" gorm:"column:user_id"`
	BorrowDate time.Time      `json:"borrow_date"`
	ReturnDate time.Time      `json:"return_date"`
	State      string         `json:"state"` //未归还：borrowing, 已归还：returned
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type BookInfo struct {
	ID       uint   `json:"id"`
	Bookname string `json:"bookname"`
	Author   string `json:"author"`
	BorrNum  int    `json:"borr_num"`
	NowNum   int    `json:"now_num"`
}

type UserInfo struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Telenum    string `json:"telenum"`
	Type       string `json:"type"` //user为普通用户，admin为管理员
	BorrRecNum int    `json:"borr_rec_num"`
}
