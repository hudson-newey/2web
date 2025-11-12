package cms

import (
	"fmt"

	"github.com/hudson-newey/2web-cli/src/shell"
)

func syncRemote(oneDrivePath string) {
	remotePath := fmt.Sprintf("OneDrive:%s", oneDrivePath)
	localPath := "./src/"

	shell.ExecuteCommand("rclone", "copy", "-P", remotePath, localPath)
}
