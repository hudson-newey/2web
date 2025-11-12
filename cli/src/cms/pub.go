package cms

import (
	"fmt"
)

// TODO: I should add an actual implementation that doesn't rely on calling
// rclone directly.
func AddCmsSource(args []string) {
	// In this example I have added a OneDrive remote called "OneDrive"
	// rclone copy -P OneDrive:Path/On/OneDrive /path/to/local/directory
	ensureCmsInitialized()

	if len(args) < 4 {
		fmt.Println("Error: Missing argument for AddCmsSource. Please provide the OneDrive path.")
		return
	}

	oneDrivePath := args[3]
	syncRemote(oneDrivePath)
}

// TODO
func ViewCmsSource() {
	fmt.Println("Not implemented.")
}

// TODO: Improve this implementation
func SyncCmsSources(args []string) {
	if len(args) < 4 {
		fmt.Println("Error: Missing argument for SyncCmsSources. Please provide the OneDrive path.")
		return
	}

	oneDrivePath := args[3]
	syncRemote(oneDrivePath)
}

// TODO
func RemoveCmsSource() {
	fmt.Println("Not implemented.")
}
