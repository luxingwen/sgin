package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	secretKey = "your-secret-key"
)

// GenerateToken 生成 JWT token
func GenerateToken(userID string) (string, error) {
	// 定义 JWT 的有效期限
	expirationTime := time.Now().Add(24 * time.Hour) // 设置为 24 小时有效期，可根据需求调整

	// 创建 token 的声明部分
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expirationTime.Unix(),
	}

	// 使用 HS256 算法进行签名
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥对 token 进行签名，生成字符串格式的 token
	tokenString, err := token.SignedString([]byte(secretKey)) // 使用与验证时相同的密钥
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
