package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

// HashPasswordWithSalt 使用盐（salt）来哈希密码
func HashPasswordWithSalt(password, salt string) string {
	h := hmac.New(sha256.New, []byte(salt))
	h.Write([]byte(password))
	hashed := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(hashed)
}

// CheckPasswordHashWithSalt 验证密码与其哈希值是否匹配
func CheckPasswordHashWithSalt(password, hashed, salt string) bool {
	expectedHash := HashPasswordWithSalt(password, salt)
	return hmac.Equal([]byte(hashed), []byte(expectedHash))
}
