package commands

import (
	"fmt"

	"osto-login-cli/internal/appctx"
	apperrors "osto-login-cli/internal/errors"
)

// HandleWhoami prints the current session's user info, or an error if the
// session has expired since it was created.
func HandleWhoami(ctx *appctx.AppContext) {
	if ctx.CurrentUser == nil || ctx.CurrentSess == nil {
		printErr(apperrors.ErrNotLoggedIn)
		return
	}
	if ctx.CurrentSess.IsExpired() {
		fmt.Println("Your session has expired. Please log in again.")
		ctx.ClearSession()
		return
	}

	u := ctx.CurrentUser
	fmt.Println("Username:   ", u.Username)
	fmt.Println("User ID:    ", u.ID)
	fmt.Println("2FA enabled:", u.TOTPEnabled)
	fmt.Println("Session expires:", ctx.CurrentSess.ExpiresAt.Format("2006-01-02 15:04:05"))
}
