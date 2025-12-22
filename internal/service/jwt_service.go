package service

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Cheerdoge/library-manage-system/internal/model"
	"github.com/golang-jwt/jwt/v5"
)

// JWTService 提供基于 HMAC-SHA256 的 JSON Web Token 签发与解析能力。
// 配置通过环境变量提供（开发可写入 go.env/.env 并由入口加载）：
//   - JWT_SECRET: 对称签名密钥（强随机、保密；生产环境不可为空）
//   - JWT_ISSUER: 签发者标识，用于校验 token 来源（如 "library-system"）
//   - JWT_EXPIRES: 令牌过期时间（单位秒，如 86400 表示 24 小时）
// 使用方式：
//   1) 登录成功后调用 GenerateToken() 生成字符串令牌返回给前端
//   2) 受保护接口在认证中间件中调用 ParseToken() 验证令牌并提取身份
//   3) 将解析出的身份映射为业务层的 Principal，以进行权限与规则判断

type Claims struct {
	// UserID: 用户唯一标识，用于后续权限与数据范围判断
	UserID uint `json:"user_id"`
	// UserName: 展示用用户名，便于日志与审计（不包含敏感信息）
	UserName string `json:"username"`
	// UserType: 用户类型（"user" 或 "admin"），用于角色判断与授权
	UserType string `json:"usertype"`
	// RegisteredClaims: 标准 JWT 字段集合，如 Issuer/IssuedAt/ExpiresAt 等
	jwt.RegisteredClaims
}

type JWTService struct {
	// secret: HMAC 对称密钥，用于签名与验签，应为强随机字节序列（建议 Base64 存储于环境变量）
	secret []byte
	// issuer: 签发者标识，与 Claims.Issuer 对应，用于校验来源
	issuer string
	// expires: 令牌有效期（Duration），超时后令牌不可用
	expires time.Duration
}

func NewJWTService() *JWTService {
	// 从环境变量加载配置；开发环境可由 godotenv.Load/Overload 预加载 go.env/.env
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "change-me-in-production"
	}
	issuer := os.Getenv("JWT_ISSUER")
	if issuer == "" {
		issuer = "library-system"
	}
	expStr := os.Getenv("JWT_EXPIRES")
	if expStr == "" {
		expStr = "86400"
	}
	expSec, _ := strconv.Atoi(expStr)
	return &JWTService{
		secret:  []byte(secret),
		issuer:  issuer,
		expires: time.Duration(expSec) * time.Second,
	}
}

// GenerateToken 签发 JWT
func (s *JWTService) GenerateToken(p *model.Principal) (string, error) {
	// 1) 构造业务 Claims：将 Principal 映射到 JWT 载荷中
	// 2) 设置标准声明：签发者（iss）、签发时间（iat）、过期时间（exp）
	// 3) 使用 HS256 与对称密钥签名，得到可传输的字符串令牌
	now := time.Now()
	claims := &Claims{
		UserID:   p.UserID,
		UserName: p.UserName,
		UserType: p.UserType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.expires)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *JWTService) ParseToken(tokenString string) (*model.Principal, error) {
	// 解析与校验流程：
	//   1) 指定我们自定义的 Claims 类型以承接载荷
	//   2) 校验签名算法必须为 HMAC（防止算法混淆攻击）
	//   3) 提供密钥并校验 Issuer；默认也会校验 exp/iat 等标准字段
	//   4) 验证通过后，将 Claims 映射回业务层 Principal
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("不支持的签名算法: %v", token.Header["alg"])
		}
		return s.secret, nil
	}, jwt.WithIssuer(s.issuer))
	if err != nil {
		// 可能的错误：签名不匹配、令牌过期、Issuer 不符、载荷格式错误等
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("无效的令牌")
	}
	// 将载荷转换为业务可用的身份对象；此对象应仅包含最小必要信息
	return &model.Principal{
		UserID:   claims.UserID,
		UserName: claims.UserName,
		UserType: claims.UserType,
	}, nil
}
