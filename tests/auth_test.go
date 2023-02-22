package tests

import (
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestName(t *testing.T) {
	password, _ := bcrypt.GenerateFromPassword([]byte("123456"), 8)
	t.Log(string(password))
}
