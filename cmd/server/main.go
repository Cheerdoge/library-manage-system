package main

import (
	"fmt"

	"github.com/Cheerdoge/library-manage-system/internal/dao"
	"github.com/Cheerdoge/library-manage-system/internal/global"
	"github.com/Cheerdoge/library-manage-system/internal/handler"
	"github.com/Cheerdoge/library-manage-system/internal/router"
	"github.com/Cheerdoge/library-manage-system/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 优先从 go.env 加载环境变量（开发环境）
	_ = godotenv.Overload("go.env")

	fmt.Println("正在初始化配置")
	global.InitConfig()

	db, err := global.InitDB()
	if err != nil {
		fmt.Printf("无法连接到数据库: %v\n", err)
		return
	}

	err = global.InitAdmin(db)
	if err != nil {
		fmt.Printf("无法初始化管理员账号: %v\n", err)
		return
	}
	fmt.Println("数据库初始化成功,管理员账号已创建,用户名:admin,密码:admin123")

	//dao层初始化
	userdao := dao.NewUserDao(db)
	bookdao := dao.NewBookDao(db)
	borrowrecorddao := dao.NewBorrowRecordDao(db)
	sessiondao := dao.NewSessionDao(db)

	//service层初始化
	bookservice := service.NewBookService(bookdao)
	sessionservice := service.NewSessionService(sessiondao)
	userservice := service.NewUserService(userdao, sessionservice)
	borrowservice := service.NewBorrowService(db, borrowrecorddao, userservice, bookservice)

	//handler层初始化
	authhandler := handler.NewAuthHandler(sessionservice, userservice)
	userhandler := handler.NewUserHandler(userservice)
	bookhandler := handler.NewBookHandler(bookservice)
	borrowhandler := handler.NewBorrowHandler(borrowservice)

	// 启动服务器
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{ // 允许的请求源
			"http://localhost:5173", // 前端vite的默认启动地址
			"http://localhost:3000", // 前端自己定义的启动地址
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},                   // 允许的请求方法
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"}, // 允许的请求头
		AllowCredentials: true,
	}))

	router.RegisterRoutes(r, bookhandler, borrowhandler, userhandler, authhandler, sessionservice)

	port := ":" + global.AppConfig.Server.Port
	fmt.Printf("服务器正在运行，监听端口 %s\n", port)
	if err := r.Run(port); err != nil {
		fmt.Printf("服务器启动失败: %v\n", err)
	}

	fmt.Println("正在与数据库断开连接")
	err = global.CloseDB(db)
	if err != nil {
		fmt.Printf("关闭数据库连接时出错: %v\n", err)
		return
	}
	fmt.Println("数据库连接已关闭")
}
