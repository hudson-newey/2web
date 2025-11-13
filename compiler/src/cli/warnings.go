package cli

import "fmt"

func PrintWarning(message string) {
	if GetArgs().IsSilent {
		return
	}

	// Add a double new line so that the warning message is emphasized in the
	// compiler logs.
	fmt.Printf("\n\033[33m[Warning]\033[0m %s\n", message)
}
