package tests

import (
	"testing"

	"osto-login-cli/internal/auth"
)

func TestHashAndVerifyPassword(t *testing.T) {
	hash, err := auth.HashPassword("correct-horse-battery-staple")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	if !auth.VerifyPassword(hash, "correct-horse-battery-staple") {
		t.Error("expected correct password to verify")
	}
	if auth.VerifyPassword(hash, "wrong-password") {
		t.Error("expected incorrect password to fail verification")
	}
}
