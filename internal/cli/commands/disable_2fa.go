package commands

import (
	"fmt"

	"github.com/chzyer/readline"

	"osto-login-cli/internal/appctx"
	"osto-login-cli/internal/auth"
	apperrors "osto-login-cli/internal/errors"
)

// HandleDisable2FA requires the user's password again (not just being
// logged in) plus a valid current TOTP code before turning 2FA off, since
// this weakens the account's security posture.
func HandleDisable2FA(ctx *appctx.AppContext, rl *readline.Instance) {
	if !ctx.IsLoggedIn() {
		printErr(apperrors.ErrNotLoggedIn)
		return
	}
	if !ctx.CurrentUser.TOTPEnabled {
		printErr(apperrors.ErrTOTPNotEnabled)
		return
	}

	password := promptPassword(rl, "Confirm your password: ")
	if !auth.VerifyPassword(ctx.CurrentUser.PasswordHash, password) {
		fmt.Println("Incorrect password.")
		return
	}

	code := prompt(rl, "Current 2FA code: ")
	if !auth.ValidateCode(ctx.CurrentUser.TOTPSecret, code) {
		fmt.Println("Invalid 2FA code.")
		return
	}

	if err := ctx.UserRepo.ClearTOTP(ctx.CurrentUser.ID); err != nil {
		printErr(err)
		return
	}
	ctx.CurrentUser.TOTPEnabled = false
	ctx.CurrentUser.TOTPSecret = ""

	fmt.Println("2FA disabled.")
}
