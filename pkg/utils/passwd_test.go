package utils

import (
	"testing"
)

func TestPasswordHashingWithSalt(t *testing.T) {
	password := "mysecretpassword"
	salt := "somesalt"

	hashed := HashPasswordWithSalt(password, salt)
	t.Log(hashed)
	if !CheckPasswordHashWithSalt(password, hashed, salt) {
		t.Errorf("Password hashing failed")
	}
}
