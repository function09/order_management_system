package auth

import (
	"testing"
	"time"
)

func TestVerifyPassword(t *testing.T) {
	hashedPassword, err := HashPassword("password")

	if err != nil {
		t.Fatalf("token generation error: %s", err)
	}

	t.Run("correct password should match", func(t *testing.T) {
		if err := VerifyPassword(hashedPassword, "password"); err != nil {
			t.Error("expected passwords to match")
		}
	})

	t.Run("wrong password should return an error", func(t *testing.T) {
		if err := VerifyPassword(hashedPassword, "wrongpassword"); err == nil {
			t.Error("expected error for wrong password")
		}
	})
}

func TestValidateToken(t *testing.T) {
	token, err := GenerateToken("user", "secret", time.Hour)

	if err != nil {
		t.Fatalf("Error generating token %s", err)
	}

	t.Run("valid token returns correct claims", func(t *testing.T) {

		claims, err := ValidateToken(token, "secret")
		want := "user"

		if err != nil {
			t.Fatalf("Error validating token %s", err)
		}

		if claims.Subject != want {
			t.Errorf("got %q want %q", claims.Subject, want)
		}
	})

	t.Run("wrong secret should return an error", func(t *testing.T) {
		_, err := ValidateToken(token, "wrongToken")

		if err == nil {
			t.Error("Expected error for wrong secret, got nil")
		}
	})

	t.Run("expired token should return an error", func(t *testing.T) {
		token, err := GenerateToken("user", "secret", -time.Hour)

		if err != nil {
			t.Fatalf("token generation error: %s", err)
		}

		_, err = ValidateToken(token, "secret")

		if err == nil {
			t.Error("Expected token expiration error, got nil")
		}
	})
}
