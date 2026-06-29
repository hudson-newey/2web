package logger

import (
	"fmt"
)

func PrintWarning(msg string) (n int, err error) {
	return fmt.Printf("\033[33m[Warning]\033[0m %s\n", msg)
}
