package utils

import (
	"testing"
)

func TestGenerateToken(t *testing.T) {
	token, err := GenerateToken("123")
	if err != nil {
		t.Error(err)
	}
	t.Log(token)
}
