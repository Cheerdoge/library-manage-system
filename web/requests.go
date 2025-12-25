package web

//===============Book==============

//新增图书请求
type AddBook struct {
	Bookname string `json:"bookname"`
	Author   string `json:"author"`
	Num      int    `json:"num"`
}

//查找图书请求
// type FindBook struct {
// 	Bookname string `json:"bookname"`
// }

// type FindBookById struct {
// 	BookId uint `json:"bookid"`
// }

//删除图书请求，分部分和全部
type DelBook struct {
	Id  uint `json:"id"`
	Num uint `json:"num"`
}

//更新图书信息，主要是数量
type UpdateBook struct {
	Id  uint `json:"id"`
	Num int  `json:"num"`
}

//================User===============

//用户注册请求
type AddUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//用户登录请求
type LoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
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

//用户注销请求
type DelUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//管理员获取指定用户信息请求
// type GetUserInfo struct {
// 	Username string `json:"username"`
// }

//================BorrowRecord===============

//新增借书记录请求
type BorrowBook struct {
	BookNum int `json:"book_num"`
}
