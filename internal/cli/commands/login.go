package commands

import (
	"fmt"

	"github.com/chzyer/readline"

	"osto-login-cli/internal/appctx"
	"osto-login-cli/internal/auth"
	apperrors "osto-login-cli/internal/errors"
)

// HandleLogin walks: password check -> lockout check -> optional 2FA ->
// session creation -> auto-display of the user summary.
func HandleLogin(ctx *appctx.AppContext, rl *readline.Instance) {
	if ctx.IsLoggedIn() {
		printErr(apperrors.ErrAlreadyLoggedIn)
		return
	}

	username := prompt(rl, "Username: ")
	password := promptPassword(rl, "Password: ")

	user, err := ctx.UserRepo.GetByUsername(username)
	if err != nil {
		// Deliberately vague: don't reveal whether the username exists.
		fmt.Println("Invalid username or password.")
		return
	}

	if auth.IsLocked(user.LockedUntil) {
		fmt.Println("Account locked. Try again later.")
		return
	}

	if !auth.VerifyPassword(user.PasswordHash, password) {
		user.FailedAttempts++
		lockedUntil := auth.ComputeLockout(user.FailedAttempts, ctx.Config.LockoutThreshold, ctx.Config.LockoutDuration)
		if err := ctx.UserRepo.UpdateFailedAttempts(user.ID, user.FailedAttempts, lockedUntil); err != nil {
			printErr(err)
			return
		}
		fmt.Println("Invalid username or password.")
		if lockedUntil != nil {
			fmt.Println("Too many failed attempts — account is now locked.")
		}
		return
	}

	if user.TOTPEnabled {
		code := prompt(rl, "2FA Code: ")
		if !auth.ValidateCode(user.TOTPSecret, code) {
			fmt.Println("Invalid 2FA code.")
			return
		}
	}

	// Success: reset attempts, create session, update last_login.
	if err := ctx.UserRepo.ResetFailedAttempts(user.ID); err != nil {
		printErr(err)
		return
	}
	if err := ctx.UserRepo.UpdateLastLogin(user.ID); err != nil {
		printErr(err)
		return
	}

	sess := auth.NewSession(user.ID, ctx.Config.SessionTimeout)
	if err := ctx.SessionRepo.Create(sess); err != nil {
		printErr(err)
		return
	}

	ctx.CurrentUser = user
	ctx.CurrentSess = &sess

	fmt.Printf("Welcome back, %s.\n", user.Username)
	printUserSummary(ctx)
}

// printUserSummary is shown automatically right after a successful login.
func printUserSummary(ctx *appctx.AppContext) {
	u := ctx.CurrentUser
	fmt.Println("---")
	fmt.Println("Username:   ", u.Username)
	fmt.Println("2FA enabled:", u.TOTPEnabled)
	if u.LastLoginAt != nil {
		fmt.Println("Last login: ", u.LastLoginAt.Format("2006-01-02 15:04:05"))
	} else {
		fmt.Println("Last login:  (first login)")
	}
	fmt.Println("Session expires:", ctx.CurrentSess.ExpiresAt.Format("2006-01-02 15:04:05"))
	fmt.Println("---")
}
