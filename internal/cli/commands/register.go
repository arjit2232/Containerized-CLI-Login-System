package commands

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"

	"github.com/chzyer/readline"

	"osto-login-cli/internal/appctx"
	"osto-login-cli/internal/auth"
	apperrors "osto-login-cli/internal/errors"
)

var usernameRe = regexp.MustCompile(`^[a-zA-Z0-9_]{3,50}$`)

// HandleRegister creates a new user account. It does not log the user in
// automatically — that keeps register/login as two clearly separate,
// independently testable flows.
func HandleRegister(ctx *appctx.AppContext, rl *readline.Instance) {
	if ctx.IsLoggedIn() {
		printErr(apperrors.ErrAlreadyLoggedIn)
		return
	}

	username := prompt(rl, "Choose a username: ")
	if !usernameRe.MatchString(username) {
		printErr(apperrors.ErrUsernameInvalid)
		return
	}

	password := promptPassword(rl, "Choose a password: ")
	if len(password) < 8 {
		printErr(apperrors.ErrPasswordTooWeak)
		return
	}
	confirm := promptPassword(rl, "Confirm password: ")
	if password != confirm {
		fmt.Println("Error: passwords do not match.")
		return
	}

	if _, err := ctx.UserRepo.GetByUsername(username); err == nil {
		printErr(apperrors.ErrUsernameTaken)
		return
	} else if !errors.Is(err, apperrors.ErrUserNotFound) && !errors.Is(err, sql.ErrNoRows) {
		printErr(err)
		return
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		printErr(err)
		return
	}

	if _, err := ctx.UserRepo.CreateUser(username, hash); err != nil {
		printErr(err)
		return
	}

	fmt.Printf("Account '%s' created. Run 'login' to sign in.\n", username)
}
