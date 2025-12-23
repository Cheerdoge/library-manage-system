package main

import (
	"fmt"

	"github.com/Cheerdoge/library-manage-system/internal/global"
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

	// 启动服务器
	fmt.Println("正在启动服务器")

	fmt.Println("正在与数据库断开连接")
	err = global.CloseDB(db)
	if err != nil {
		fmt.Printf("关闭数据库连接时出错: %v\n", err)
		return
	}
	fmt.Println("数据库连接已关闭")
}
