package cli

import (
	"fmt"
	"os"
)

func PrintError(code int, message string) {
	fmt.Printf("\033[31m[Error]\033[0m: %s\n", message)
	os.Exit(code)
}
