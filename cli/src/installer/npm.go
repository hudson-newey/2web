package installer

import (
	"fmt"
	"os/exec"

	"github.com/hudson-newey/2web/_shared/logger"
)

func installNpmPackage(name string) {
	cmd := exec.Command("npm", "install", name)
	stdout, err := cmd.Output()

	if err != nil {
		errorMsg := fmt.Sprintf("failed to install package '%s': %s", name, err)
		logger.PrintError(errorMsg)
	}

	logger.Println(string(stdout))
}
