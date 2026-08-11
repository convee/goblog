package tests

import (
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestGeneratedPasswordCanBeVerified(t *testing.T) {
	password, err := bcrypt.GenerateFromPassword([]byte("a-test-password"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	if err := bcrypt.CompareHashAndPassword(password, []byte("a-test-password")); err != nil {
		t.Fatalf("generated password hash cannot be verified: %v", err)
	}
}
