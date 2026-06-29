package logger

import (
	"fmt"
	"os"
)

func PrintError(msg string) {
	fmt.Fprintf(os.Stderr, "\033[31m[Error]\033[0m: %s", msg)
	os.Exit(1)
}
