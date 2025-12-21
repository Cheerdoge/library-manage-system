package model

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	Id        uint   `json:"id"`
	Bookname  string `json:"bookname"`
	Author    string `json:"author"`
	Num       int    `json:"num"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type User struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	Password  string `json:"-"`
	Telenum   string `json:"telenum"`
	Type      uint   `json:"type"` //1为普通用户，0为管理员
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type BorrowRecord struct {
	Id         uint `json:"id"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	Bookid     uint           `json:"book_id"`
	Userid     uint           `json:"user_id"`
	BorrowData time.Time      `json:"borrow_data"`
	ReturnData time.Time      `json:"return_data"`
	State      string         `json:"state"` //未归还：borrowing, 已归还：returned
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
