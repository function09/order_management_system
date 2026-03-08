package auth

import (
	"testing"
)

func TestPasswordHash(t *testing.T) {
	hash, err := HashPassword("password")

	if err != nil {
		t.Fatalf("expected no error got %v", err)
	}

	if hash == "password" {
		t.Error("hash should not equal original password")
	}

	if hash == "" {
		t.Error("hash should not be empty")
	}
}

func TestVerifyPassword(t *testing.T) {
	hashedPassword, _ := HashPassword("password")

	if err := VerifyPassword(hashedPassword, "password"); err != nil {
		t.Error("expected passwords to match")
	}

}
