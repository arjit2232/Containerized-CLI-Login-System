package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	apperrors "osto-login-cli/internal/errors"
	"osto-login-cli/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser inserts a new user and returns the generated ID.
func (r *UserRepository) CreateUser(username, passwordHash string) (string, error) {
	id := uuid.NewString()
	_, err := r.db.Exec(
		`INSERT INTO users (id, username, password_hash) VALUES (?, ?, ?)`,
		id, username, passwordHash,
	)
	if err != nil {
		return "", err
	}
	return id, nil
}

// GetByUsername fetches a user by username.
func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	row := r.db.QueryRow(
		`SELECT id, username, password_hash, totp_secret, totp_enabled,
		        failed_attempts, locked_until, created_at, last_login_at
		 FROM users WHERE username = ?`,
		username,
	)
	return scanUser(row)
}

// GetByID fetches a user by ID.
func (r *UserRepository) GetByID(id string) (*models.User, error) {
	row := r.db.QueryRow(
		`SELECT id, username, password_hash, totp_secret, totp_enabled,
		        failed_attempts, locked_until, created_at, last_login_at
		 FROM users WHERE id = ?`,
		id,
	)
	return scanUser(row)
}

// UpdateFailedAttempts persists the failed-login counter and lockout expiry.
func (r *UserRepository) UpdateFailedAttempts(id string, attempts int, lockedUntil *time.Time) error {
	_, err := r.db.Exec(
		`UPDATE users SET failed_attempts = ?, locked_until = ? WHERE id = ?`,
		attempts, lockedUntil, id,
	)
	return err
}

// ResetFailedAttempts clears the failed-login counter and lockout on success.
func (r *UserRepository) ResetFailedAttempts(id string) error {
	_, err := r.db.Exec(
		`UPDATE users SET failed_attempts = 0, locked_until = NULL WHERE id = ?`,
		id,
	)
	return err
}

// UpdateLastLogin stamps the current time as the last successful login.
func (r *UserRepository) UpdateLastLogin(id string) error {
	_, err := r.db.Exec(
		`UPDATE users SET last_login_at = ? WHERE id = ?`,
		time.Now(), id,
	)
	return err
}

// SetTOTPSecret stores a (not-yet-confirmed) TOTP secret for a user.
func (r *UserRepository) SetTOTPSecret(id, secret string) error {
	_, err := r.db.Exec(
		`UPDATE users SET totp_secret = ? WHERE id = ?`,
		secret, id,
	)
	return err
}

// SetTOTPEnabled flips the totp_enabled flag.
func (r *UserRepository) SetTOTPEnabled(id string, enabled bool) error {
	_, err := r.db.Exec(
		`UPDATE users SET totp_enabled = ? WHERE id = ?`,
		enabled, id,
	)
	return err
}

// ClearTOTP disables 2FA and wipes the stored secret.
func (r *UserRepository) ClearTOTP(id string) error {
	_, err := r.db.Exec(
		`UPDATE users SET totp_enabled = FALSE, totp_secret = NULL WHERE id = ?`,
		id,
	)
	return err
}

func scanUser(row *sql.Row) (*models.User, error) {
	var u models.User
	var totpSecret sql.NullString
	var lockedUntil sql.NullTime
	var lastLoginAt sql.NullTime

	err := row.Scan(
		&u.ID, &u.Username, &u.PasswordHash, &totpSecret, &u.TOTPEnabled,
		&u.FailedAttempts, &lockedUntil, &u.CreatedAt, &lastLoginAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperrors.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	u.TOTPSecret = totpSecret.String
	if lockedUntil.Valid {
		t := lockedUntil.Time
		u.LockedUntil = &t
	}
	if lastLoginAt.Valid {
		t := lastLoginAt.Time
		u.LastLoginAt = &t
	}

	return &u, nil
}
