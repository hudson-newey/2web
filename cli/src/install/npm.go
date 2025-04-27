package install

import (
	"fmt"
	"os/exec"
)

func installNpmPackage(name string) {
	cmd := exec.Command("npm", "install", name)
	stdout, err := cmd.Output()

	if err != nil {
		errorMsg := fmt.Errorf("failed to install package '%s': %s", name, err)
		panic(errorMsg)
	}

	fmt.Println(string(stdout))
}
