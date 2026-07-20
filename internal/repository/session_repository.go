package repository

import (
	"database/sql"
	"errors"
	"time"

	apperrors "osto-login-cli/internal/errors"
	"osto-login-cli/internal/models"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// Create persists a new session row.
func (r *SessionRepository) Create(s models.Session) error {
	_, err := r.db.Exec(
		`INSERT INTO sessions (id, user_id, token, created_at, expires_at)
		 VALUES (?, ?, ?, ?, ?)`,
		s.ID, s.UserID, s.Token, s.CreatedAt, s.ExpiresAt,
	)
	return err
}

// GetActiveSession returns the non-revoked, non-expired session for a token.
func (r *SessionRepository) GetActiveSession(token string) (*models.Session, error) {
	row := r.db.QueryRow(
		`SELECT id, user_id, token, created_at, expires_at, revoked_at
		 FROM sessions
		 WHERE token = ? AND revoked_at IS NULL AND expires_at > ?`,
		token, time.Now(),
	)

	var s models.Session
	var revokedAt sql.NullTime
	err := row.Scan(&s.ID, &s.UserID, &s.Token, &s.CreatedAt, &s.ExpiresAt, &revokedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperrors.ErrSessionNotFound
	}
	if err != nil {
		return nil, err
	}
	if revokedAt.Valid {
		t := revokedAt.Time
		s.RevokedAt = &t
	}
	return &s, nil
}

// RevokeSession marks a session as revoked (used on logout).
func (r *SessionRepository) RevokeSession(id string) error {
	_, err := r.db.Exec(
		`UPDATE sessions SET revoked_at = ? WHERE id = ? AND revoked_at IS NULL`,
		time.Now(), id,
	)
	return err
}

// RevokeAllForUser revokes every active session belonging to a user
// (handy for "log out everywhere" or account-lock scenarios).
func (r *SessionRepository) RevokeAllForUser(userID string) error {
	_, err := r.db.Exec(
		`UPDATE sessions SET revoked_at = ? WHERE user_id = ? AND revoked_at IS NULL`,
		time.Now(), userID,
	)
	return err
}