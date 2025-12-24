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
	UserName   string `json:"username"`
	Password   string `json:"-"`
	Telenum    string `json:"telenum"`
	NowBorrNum int    `json:"now_borr_num"`
	OverdueNum int    `json:"overdue_num"`
	IsAdmin    bool   `json:"is_admin"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

type BorrowRecord struct {
	ID           uint `json:"id"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	BookID       uint           `json:"book_id" gorm:"column:book_id"`
	UserID       uint           `json:"user_id" gorm:"column:user_id"`
	BookNum      int            `json:"book_num"`
	BorrowDate   time.Time      `json:"borrow_date"`
	ReturnDate   time.Time      `json:"return_date"`
	ShouldReturn time.Time      `json:"should_return"`
	State        string         `json:"state"` //未归还：borrowing, 已归还：returned
}

type BookInfo struct {
	ID       uint   `json:"id"`
	Bookname string `json:"bookname"`
	Author   string `json:"author"`
	SumNum   int    `json:"sum_num"`
	BorrNum  int    `json:"borr_num"`
	NowNum   int    `json:"now_num"`
}

type UserInfo struct {
	ID         uint   `json:"id"`
	UserName   string `json:"username"`
	Telenum    string `json:"telenum"`
	IsAdmin    bool   `json:"is_admin"`
	NowBorrNum int    `json:"now_borr_num"`
	OverdueNum int    `json:"overdue_num"`
}

type Session struct {
	ID        uint      `json:"id"`
	UserName  string    `json:"username"`
	UserID    uint      `json:"user_id"`
	IsAdmin   bool      `json:"is_admin"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
