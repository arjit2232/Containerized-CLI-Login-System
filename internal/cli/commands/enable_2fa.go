package commands

import (
	"fmt"

	"github.com/chzyer/readline"

	"osto-login-cli/internal/appctx"
	"osto-login-cli/internal/auth"
	apperrors "osto-login-cli/internal/errors"
)

// HandleEnable2FA generates a TOTP secret, shows it (as text + otpauth URI
// for QR generation) and requires the user to prove they've set it up by
// entering a valid code before it's actually turned on.
func HandleEnable2FA(ctx *appctx.AppContext, rl *readline.Instance) {
	if !ctx.IsLoggedIn() {
		printErr(apperrors.ErrNotLoggedIn)
		return
	}
	if ctx.CurrentUser.TOTPEnabled {
		printErr(apperrors.ErrTOTPAlreadyEnabled)
		return
	}

	secret, qrURI, err := auth.GenerateSecret(ctx.CurrentUser.Username)
	if err != nil {
		printErr(err)
		return
	}

	if err := ctx.UserRepo.SetTOTPSecret(ctx.CurrentUser.ID, secret); err != nil {
		printErr(err)
		return
	}

	fmt.Println("Scan this URI with your authenticator app (or paste into a QR generator):")
	fmt.Println(qrURI)
	fmt.Println("Or enter the secret manually:", secret)

	code := prompt(rl, "Enter the 6-digit code to confirm: ")
	if !auth.ValidateCode(secret, code) {
		fmt.Println("Invalid code. 2FA was not enabled — run 'enable-2fa' again to retry.")
		return
	}

	if err := ctx.UserRepo.SetTOTPEnabled(ctx.CurrentUser.ID, true); err != nil {
		printErr(err)
		return
	}
	ctx.CurrentUser.TOTPEnabled = true

	fmt.Println("2FA enabled successfully.")
}
