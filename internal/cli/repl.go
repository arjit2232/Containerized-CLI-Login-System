package cli

import (
	"fmt"
	"strings"

	"github.com/chzyer/readline"

	"osto-login-cli/internal/appctx"
	"osto-login-cli/internal/cli/commands"
)

// Run starts the interactive REPL loop until the user exits (Ctrl+D,
// Ctrl+C, or the "exit" command).
func Run(ctx *appctx.AppContext) error {
	rl, err := readline.New("> ")
	if err != nil {
		return err
	}
	defer rl.Close()

	fmt.Println("osto-login-cli — type 'help' for a list of commands")

	for {
		line, err := rl.Readline()
		if err != nil { // Ctrl+D / Ctrl+C
			break
		}
		cmd := strings.TrimSpace(line)
		if cmd == "" {
			continue
		}
		if dispatch(ctx, cmd, rl) {
			break // "exit" was called
		}
	}
	return nil
}

// dispatch routes a command to its handler. It returns true if the REPL
// loop should terminate (the "exit" command).
func dispatch(ctx *appctx.AppContext, cmd string, rl *readline.Instance) (exit bool) {
	switch cmd {
	case "help":
		commands.HandleHelp(ctx)
	case "exit":
		fmt.Println("Goodbye.")
		return true
	case "logout":
		commands.HandleLogout(ctx)
	case "register":
		commands.HandleRegister(ctx, rl)
	case "login":
		commands.HandleLogin(ctx, rl)
	case "whoami":
		commands.HandleWhoami(ctx)
	case "enable-2fa":
		commands.HandleEnable2FA(ctx, rl)
	case "disable-2fa":
		commands.HandleDisable2FA(ctx, rl)
	default:
		fmt.Println("Unknown command. Type 'help' for a list of commands.")
	}
	return false
}
