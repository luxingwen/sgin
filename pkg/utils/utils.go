package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"strconv"
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

// ParseToken 解析 JWT token
func ParseToken(tokenString string) (jwt.MapClaims, error) {
	// 解析 token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil // 使用与生成 token 时相同的密钥
	})
	if err != nil {
		return nil, err
	}

	// 获取 token 中的声明部分
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err
	}

	return claims, nil
}

// 解析token返回user_id
func ParseTokenGetUserID(tokenString string) (string, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", err
	}

	return userID, nil
}

// 生成验证码
func GenerateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(999999)
	return strconv.Itoa(code)
}

// SignBody 签名
func SignBody(body, secretKey []byte) string {
	mac := hmac.New(sha256.New, secretKey)
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}
