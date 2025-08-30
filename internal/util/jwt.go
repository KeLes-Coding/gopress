// package util 提供项目需要的各种工具函数，例如 JWT 处理、数据加密等。
package util

import (
	"errors"
	"time"

	"github.com/KeLes-Coding/gopress/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

// MyClaims 是一个自定义的 Claims 结构体，用于在 JWT 中存放我们自己的业务数据。
// 我们选择嵌入 jwt.RegisteredClaims 来包含 JWT 规范中预定义的字段（如 iss, exp, sub 等）。
type MyClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 函数用于根据用户 ID 和用户名生成一个新的 JWT。
func GenerateToken(userID uint, username string) (string, error) {
	// 创建自定义的 claims
	claims := MyClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			// 设置过期时间，例如 7 天后过期
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			// 设置签发时间
			IssuedAt: jwt.NewNumericDate(time.Now()),
			// 设置签发人
			Issuer: "gopress",
		},
	}

	// 使用 HS256 签名方法创建一个新的 token 实例
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用配置文件中定义的密钥来为 token 签名，并获取完整的 token 字符串
	signedToken, err := token.SignedString([]byte(config.Conf.Server.JWTSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// ParseToken 函数用于解析和验证一个 JWT 字符串。
// 如果 token 有效，它会返回包含用户信息的 MyClaims 指针。
func ParseToken(tokenString string) (*MyClaims, error) {
	// 使用密钥和自定义的 Claims 结构体解析 token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Conf.Server.JWTSecret), nil
	})

	if err != nil {
		// 如果解析过程中发生错误，则返回错误
		return nil, err
	}

	// 检查 token 是否有效，并从中断言出我们自定义的 claims
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
