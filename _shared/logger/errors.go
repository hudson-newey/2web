package logger

import (
	"fmt"
	"os"
)

func PrintError(msg string) {
	fmt.Printf("\033[31m[Error]\033[0m: %s\n", msg)
	os.Exit(1)
}
