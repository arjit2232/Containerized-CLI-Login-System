package commands

import (
	"fmt"

	"osto-login-cli/internal/appctx"
	apperrors "osto-login-cli/internal/errors"
)

// HandleLogout revokes the current session in the DB and clears the
// in-memory logged-in state.
func HandleLogout(ctx *appctx.AppContext) {
	if !ctx.IsLoggedIn() {
		printErr(apperrors.ErrNotLoggedIn)
		return
	}

	if err := ctx.SessionRepo.RevokeSession(ctx.CurrentSess.ID); err != nil {
		printErr(err)
		return
	}

	fmt.Printf("Goodbye, %s.\n", ctx.CurrentUser.Username)
	ctx.ClearSession()
}
