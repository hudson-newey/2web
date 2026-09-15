package logger

import (
	"fmt"
	"os"
)

func PrintError(msg string) {
	// The trailing newline matters: without it, shells with partial line
	// markers (e.g. zsh) print a stray "%" after every failed command.
	fmt.Fprintf(os.Stderr, "\n\033[31m[Error]\033[0m: %s\n", msg)
	os.Exit(1)
}
