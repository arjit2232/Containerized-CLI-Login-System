package tests

import (
	"testing"
	"time"

	"osto-login-cli/internal/auth"
)

func TestComputeLockout(t *testing.T) {
	const threshold = 5
	const duration = 15 * time.Minute

	if got := auth.ComputeLockout(4, threshold, duration); got != nil {
		t.Errorf("expected no lockout below threshold, got %v", got)
	}

	got := auth.ComputeLockout(5, threshold, duration)
	if got == nil {
		t.Fatal("expected a lockout time once threshold is reached")
	}
	if !got.After(time.Now()) {
		t.Error("expected lockedUntil to be in the future")
	}
}

func TestIsLocked(t *testing.T) {
	future := time.Now().Add(10 * time.Minute)
	past := time.Now().Add(-10 * time.Minute)

	if auth.IsLocked(nil) {
		t.Error("nil lockedUntil should mean not locked")
	}
	if !auth.IsLocked(&future) {
		t.Error("future lockedUntil should mean locked")
	}
	if auth.IsLocked(&past) {
		t.Error("past lockedUntil should mean not locked")
	}
}
