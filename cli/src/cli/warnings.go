package cli

import "fmt"

func PrintWarning(message string) {
	fmt.Printf("\033[33m[Warning]\033[0m %s\n", message)
}
