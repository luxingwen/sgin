package utils

import (
	"crypto/hmac"
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math/rand"
	"strconv"
	"time"

	"sgin/pkg/config"

	jwt "github.com/golang-jwt/jwt/v4"
)

// 从配置加载 JWT 秘钥，若未配置则回退到 PasswdKey
func jwtSecret() []byte {
	cfg := config.GetConfig()
	if cfg != nil && cfg.PasswdKey != "" {
		return []byte(cfg.PasswdKey)
	}
	// 最小化改动：保持兼容旧行为
	return []byte("your-secret-key")
}

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
	tokenString, err := token.SignedString(jwtSecret()) // 使用与验证时相同的密钥
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken 解析 JWT token
func ParseToken(tokenString string) (jwt.MapClaims, error) {
	// 解析 token 并校验签名算法与有效性
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 仅接受 HMAC 系列算法，防止 none 算法攻击
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret(), nil
	})
	if err != nil {
		return nil, err
	}

	// 获取 token 中的声明部分并确保 token.Valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
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
	// 使用加密安全的随机数，生成 6 位数字码
	// 如果失败，回退到非安全随机（极少发生）
	const n = 6
	max := 1000000
	var num int
	var buf [4]byte
	if _, err := crand.Read(buf[:]); err == nil {
		// 将 4 字节转为整数并取模
		num = int((uint32(buf[0])<<24 | uint32(buf[1])<<16 | uint32(buf[2])<<8 | uint32(buf[3])) % uint32(max))
	} else {
		rand.Seed(time.Now().UnixNano())
		num = rand.Intn(max)
	}
	s := strconv.Itoa(num)
	// 左侧补零到 6 位
	for len(s) < n {
		s = "0" + s
	}
	return s
}

// SignBody 签名
func SignBody(body, secretKey []byte) string {
	mac := hmac.New(sha256.New, secretKey)
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}
