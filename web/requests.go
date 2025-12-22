package web

//===============Book==============

//新增图书请求
type AddBook struct {
	Bookname string `json:"bookname"`
	Author   string `json:"author"`
	Num      int    `json:"num"`
}

//查找图书请求，逻辑：书名-->id-->书
type FindBook struct {
	Bookname string `json:"bookname"`
}

//删除图书请求，分部分和全部
type DelBook struct {
	Id  uint `json:"id"`
	Num uint `json:"num"`
}

//更新图书信息，主要是数量
type UpdateBook struct {
	Id  uint `json:"id"`
	Num uint `json:"num"`
}

//================User===============

//用户注册请求
type AddUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Type     uint   `json:"type"` //1为普通用户，0为管理员
}

//用户登录请求
type LoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//用户登出请求
type LogoutUser struct {
	Username string `json:"username"`
}

//用户修改密码请求
type ChangePassword struct {
	Oldpassword string `json:"oldpassword"`
	Newpassword string `json:"newpassword"`
}

//用户修改信息请求
type ChangeUserInfo struct {
	Username string `json:"username"`
	Telenum  string `json:"telenum"`
}

//查找用户请求
type FindUser struct {
	Username string `json:"username"`
}

//用户注销请求
type DelUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//================BorrowRecord===============

//新增借书记录请求
type BorrowBook struct {
	Bookid uint `json:"book_id"`
	Userid uint `json:"user_id"`
}

//归还图书请求
type ReturnBook struct {
	Recordid uint `json:"record_id"`
}
