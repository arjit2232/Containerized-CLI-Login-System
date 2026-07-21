package models

import "time"

// Session mirrors the `sessions` table.
type Session struct {
	ID        string
	UserID    string
	Token     string
	CreatedAt time.Time
	ExpiresAt time.Time
	RevokedAt *time.Time
}

// IsExpired reports whether the session has passed its expiry time.
func (s Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// IsRevoked reports whether the session has been explicitly revoked.
func (s Session) IsRevoked() bool {
	return s.RevokedAt != nil
}

// IsValid reports whether the session can still be used to authenticate.
func (s Session) IsValid() bool {
	return !s.IsExpired() && !s.IsRevoked()
}
