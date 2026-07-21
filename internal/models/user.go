package models

import "time"

// User mirrors the `users` table.
type User struct {
	ID             string
	Username       string
	PasswordHash   string
	TOTPSecret     string
	TOTPEnabled    bool
	FailedAttempts int
	LockedUntil    *time.Time
	CreatedAt      time.Time
	LastLoginAt    *time.Time
}
