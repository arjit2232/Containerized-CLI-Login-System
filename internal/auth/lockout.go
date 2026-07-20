package auth

import "time"

const (
	DefaultMaxFailedAttempts = 5
	DefaultLockoutDuration   = 15 * time.Minute
)

// IsLocked reports whether the account is currently within its lockout window.
func IsLocked(lockedUntil *time.Time) bool {
	return lockedUntil != nil && time.Now().Before(*lockedUntil)
}

// ComputeLockout returns a new lockedUntil pointer once failedAttempts hits
// the configured threshold, or nil if the account should remain unlocked.
func ComputeLockout(failedAttempts, threshold int, duration time.Duration) *time.Time {
	if failedAttempts >= threshold {
		t := time.Now().Add(duration)
		return &t
	}
	return nil
}