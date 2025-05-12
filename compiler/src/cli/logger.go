package cli

import (
	"fmt"
	"log"
)

func PrintWarning(message string) {
	// add a double new line (because this is also inside a println) so that the
	// warning message is emphasized in the compiler logs
	fmt.Println()
	fmt.Println("\033[33m[Warning]\033[0m", message)
	fmt.Println()
}

func PrintBuildLog(message string) {
	if !*GetArgs().IsSilent {
		log.Println(message)
	}
}
