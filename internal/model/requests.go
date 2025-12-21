package model

//===============Book==============

//新增图书请求
type AddBook struct {
	Bookname string `json:"bookname"`
	Author   string `json:"author"`
	Num      int    `json:"num"`
}

//查找图书请求，预计逻辑：书名-->id-->书
type FindBook struct {
	Id uint `json:"id"`
}

//删除图书请求，分部分和全部
type DelBook struct {
	Id  uint `json:"id"`
	Num uint `json:"num"`
}

//================User===============

//用户注册请求
type AddUser struct {
	Username string `json:"username"`
	Password string `json:"-"`
	Type     uint   `json:"type"` //1为普通用户，0为管理员
}

//用户登录请求
type LoginUser struct {
	Username string `json:"username"`
	Password string `json:"-"`
}

//用户登出请求
type LogoutUser struct {
	Username string `json:"username"`
}

//用户修改密码请求
type ChangePassword struct {
	Username    string `json:"username"`
	Oldpassword string `json:"oldpassword"`
	Newpassword string `json:"newpassword"`
}

//用户修改信息请求
type ChangeUserInfo struct {
	Username string `json:"username"`
	Telenum  string `json:"telenum"`
}

//================BorrowRecord===============

//新增借书记录请求
type AddBorrowRecord struct {
	Bookid uint `json:"book_id"`
	Userid uint `json:"user_id"`
}

//归还图书请求
type ReturnBook struct {
	Bookid uint   `json:"book_id"`
	Userid uint   `json:"user_id"`
	Type   string `json:"type"` //归还类型，分为正常归还和逾期归还，normal和overdue
}
