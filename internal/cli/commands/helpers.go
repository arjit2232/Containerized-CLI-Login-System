package commands

import (
	"fmt"
	"strings"

	"github.com/chzyer/readline"
)

// prompt reads a single line of plain-text input.
func prompt(rl *readline.Instance, label string) string {
	rl.SetPrompt(label)
	line, err := rl.Readline()
	rl.SetPrompt("> ")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(line)
}

// promptPassword reads a line with the input masked (no echo).
func promptPassword(rl *readline.Instance, label string) string {
	pw, err := rl.ReadPassword(label)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(pw))
}

func printErr(err error) {
	fmt.Println("Error:", err)
}
