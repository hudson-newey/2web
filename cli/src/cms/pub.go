package cms

import (
	"fmt"

	"github.com/hudson-newey/2web-cli/src/shell"
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
	remotePath := fmt.Sprintf("OneDrive:%s", oneDrivePath)
	localPath := "./src/"

	shell.ExecuteCommand("rclone", "copy", "-P", remotePath, localPath)
}

// TODO
func ViewCmsSource() {
	fmt.Println("Not implemented.")
}

// TODO
func SyncCmsSources() {
	fmt.Println("Not implemented.")
}

// TODO
func RemoveCmsSource() {
	fmt.Println("Not implemented.")
}
