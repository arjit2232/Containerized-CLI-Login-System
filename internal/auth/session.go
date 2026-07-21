package auth

import (
	"time"

	"github.com/google/uuid"

	"osto-login-cli/internal/models"
)

// NewSession builds a fresh session for a user, honoring the configured
// session timeout.
func NewSession(userID string, timeout time.Duration) models.Session {
	now := time.Now()
	return models.Session{
		ID:        uuid.NewString(),
		UserID:    userID,
		Token:     uuid.NewString(),
		CreatedAt: now,
		ExpiresAt: now.Add(timeout),
	}
}
