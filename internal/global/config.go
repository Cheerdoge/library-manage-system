package global

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config 全局配置结构体
type Config struct {
	DB struct {
		User string
		Pass string
		Host string
		Port string
		Name string
	}
	Server struct {
		Port string
		Env  string // "dev" 或 "prod"
	}
	JWT struct {
		Secret  string
		Issuer  string
		Expires int // 秒数
	}
}

// AppConfig 全局配置实例
var AppConfig *Config

// InitConfig 初始化配置（从环境变量加载）
func InitConfig() {
	// 加载 .env 和 go.env，文件值覆盖系统值
	_ = godotenv.Overload(".env", "go.env")

	cfg := &Config{}

	// 数据库配置
	cfg.DB.User = getEnv("DB_USER", "root")
	cfg.DB.Pass = getEnv("DB_PASS", "")
	cfg.DB.Host = getEnv("DB_HOST", "127.0.0.1")
	cfg.DB.Port = getEnv("DB_PORT", "3306")
	cfg.DB.Name = getEnv("DB_NAME", "library_db")

	// 服务器配置
	cfg.Server.Port = getEnv("SERVER_PORT", "8080")
	cfg.Server.Env = getEnv("APP_ENV", "dev")

	// JWT 配置
	cfg.JWT.Secret = getEnv("JWT_SECRET", "change-me-in-production")
	cfg.JWT.Issuer = getEnv("JWT_ISSUER", "library-system")
	cfg.JWT.Expires = parseInt(getEnv("JWT_EXPIRES", "86400"))

	AppConfig = cfg

	// 日志输出（不含密码）
	fmt.Printf("[配置] 数据库=%s@%s:%s/%s\n",
		cfg.DB.User, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name)
	fmt.Printf("[配置] 服务器端口=%s 环境=%s\n",
		cfg.Server.Port, cfg.Server.Env)
}

// getEnv 获取环境变量，不存在时返回默认值
func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

// parseInt 将字符串转为整数，失败时返回 0
func parseInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
