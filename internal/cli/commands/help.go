package commands

import (
	"fmt"

	"osto-login-cli/internal/appctx"
)

// HandleHelp prints the command reference, adjusted for whether the user
// is currently logged in.
func HandleHelp(ctx *appctx.AppContext) {
	fmt.Println("Available commands:")
	if !ctx.IsLoggedIn() {
		fmt.Println("  register      create a new account")
		fmt.Println("  login         sign in to an existing account")
	} else {
		fmt.Println("  whoami        show info about the current session")
		fmt.Println("  enable-2fa    turn on 2FA for this account")
		fmt.Println("  disable-2fa   turn off 2FA for this account")
		fmt.Println("  logout        sign out")
	}
	fmt.Println("  help          show this message")
	fmt.Println("  exit          quit the CLI")
}
