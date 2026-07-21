// Package appctx holds AppContext, the struct threaded through every
// command handler. It lives in its own package (separate from
// internal/cli) because both internal/cli (the REPL/dispatcher) and
// internal/cli/commands (the individual handlers) need to reference it,
// and putting it in either of those would create an import cycle.
package appctx

import (
	"osto-login-cli/internal/config"
	"osto-login-cli/internal/models"
	"osto-login-cli/internal/repository"
)

// AppContext holds shared dependencies plus the current logged-in
// user/session, which command handlers mutate as the user logs in/out.
type AppContext struct {
	UserRepo    *repository.UserRepository
	SessionRepo *repository.SessionRepository
	Config      *config.Config

	CurrentUser *models.User
	CurrentSess *models.Session
}

// IsLoggedIn reports whether there's a valid, non-expired session active.
func (c *AppContext) IsLoggedIn() bool {
	return c.CurrentUser != nil && c.CurrentSess != nil && c.CurrentSess.IsValid()
}

// ClearSession drops the in-memory logged-in state (used by logout and by
// the REPL when it notices a session has expired).
func (c *AppContext) ClearSession() {
	c.CurrentUser = nil
	c.CurrentSess = nil
}
