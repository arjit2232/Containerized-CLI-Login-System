package tests

import (
	"testing"
	"time"

	"github.com/pquerna/otp/totp"

	"osto-login-cli/internal/auth"
)

func TestGenerateSecretAndValidateCode(t *testing.T) {
	secret, uri, err := auth.GenerateSecret("alice")
	if err != nil {
		t.Fatalf("GenerateSecret returned error: %v", err)
	}
	if secret == "" {
		t.Fatal("expected non-empty secret")
	}
	if uri == "" {
		t.Fatal("expected non-empty otpauth URI")
	}

	code, err := totp.GenerateCode(secret, time.Now())
	if err != nil {
		t.Fatalf("failed to generate a test code: %v", err)
	}

	if !auth.ValidateCode(secret, code) {
		t.Error("expected freshly generated code to validate")
	}
	if auth.ValidateCode(secret, "000000") {
		t.Error("expected an arbitrary code to fail validation (extremely unlikely to pass by chance)")
	}
}
