package cms

import (
	"fmt"
	"os"
	"strings"
)

// TODO: Add an actual implementation
func hasInitCmsSource() bool {
	return true
}

// TODO
func initCmsSource() {
	fmt.Println("Not implemented.")
}

func ensureCmsInitialized() {
	if !hasInitCmsSource() {
		fmt.Println("CMS source not initialized. Would you like to initialize it now? (y/n)")

		var response string
		fmt.Scanln(&response)

		if strings.ToLower(response) == "y" {
			initCmsSource()
		} else {
			fmt.Println("CMS source initialization canceled.")
			os.Exit(0)
		}

		return
	}
}
