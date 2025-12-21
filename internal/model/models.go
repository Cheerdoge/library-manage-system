package model

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	ID        uint   `json:"id"`
	Bookname  string `json:"bookname"`
	Author    string `json:"author"`
	Num       int    `json:"num"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type User struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Password  string `json:"-"`
	Telenum   string `json:"telenum"`
	Type      uint   `json:"type"` //1为普通用户，0为管理员
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
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
