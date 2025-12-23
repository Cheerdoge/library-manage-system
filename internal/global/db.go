package global

import (
	"fmt"
	"os"
	"time"

	"github.com/Cheerdoge/library-manage-system/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	// 为缺省环境变量提供安全的默认值，避免出现 :0 的端口错误
	if user == "" {
		user = "root"
	}
	if host == "" {
		host = "127.0.0.1"
	}
	if port == "" {
		port = "3306"
	}
	if name == "" {
		name = "library_db"
	}

	dsn := user + ":" + pass + "@tcp(" + host + ":" + port + ")/" + name + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(
		&model.User{},
		&model.Book{},
		&model.BorrowRecord{},
	); err != nil {
		return nil, err
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(20)                 // 最大并发连接数
	sqlDB.SetMaxIdleConns(10)                 // 最大空闲连接数
	sqlDB.SetConnMaxLifetime(5 * time.Minute) // 单连接最长存活
	sqlDB.SetConnMaxIdleTime(2 * time.Minute) // 单连接最长空闲

	return db, nil
}

func InitAdmin(db *gorm.DB) error {
	var user model.User
	result := db.FirstOrCreate(&user, "type = ?", "admin")
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected > 0 {
		fmt.Println("正在新建管理员账号")
		user.Name = "admin"
		user.Password = "admin123"
		user.Type = "admin"
		if err := db.Save(&user).Error; err != nil {
			return err
		}
	}
	return nil
}

// 关闭数据库连接
func CloseDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	if err := sqlDB.Close(); err != nil {
		return err
	}
	return nil
}
