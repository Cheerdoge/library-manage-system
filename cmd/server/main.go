package main

import (
	"fmt"

	"github.com/Cheerdoge/library-manage-system/internal/global"
	"github.com/joho/godotenv"
)

func main() {
	// 优先从 go.env 加载环境变量（开发环境）
	_ = godotenv.Load("go.env")

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

	fmt.Println("正在与数据库断开连接")
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Printf("无法获取数据库连接: %v\n", err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		fmt.Printf("关闭数据库连接时出错: %v\n", err)
		return
	}
	fmt.Println("数据库连接已关闭")
}
