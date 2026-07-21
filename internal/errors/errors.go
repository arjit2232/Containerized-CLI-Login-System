package errors

import "errors"

var (
	// Auth errors
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrAccountLocked      = errors.New("account is locked, try again later")
	ErrInvalid2FACode     = errors.New("invalid 2fa code")
	ErrNotLoggedIn        = errors.New("you are not logged in")
	ErrAlreadyLoggedIn    = errors.New("you are already logged in, run 'logout' first")

	// Registration errors
	ErrUsernameTaken   = errors.New("username is already taken")
	ErrUsernameInvalid = errors.New("username must be 3-50 characters, alphanumeric or underscore")
	ErrPasswordTooWeak = errors.New("password must be at least 8 characters")

	// 2FA errors
	ErrTOTPAlreadyEnabled = errors.New("2fa is already enabled on this account")
	ErrTOTPNotEnabled     = errors.New("2fa is not enabled on this account")

	// Repository / infra errors
	ErrUserNotFound    = errors.New("user not found")
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session has expired")
)
