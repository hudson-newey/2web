package install

import (
	"fmt"
	"os/exec"
)

func installYarnPackage(name string) {
	cmd := exec.Command("yarn", "add", name)
	stdout, err := cmd.Output()

	if err != nil {
		errorMsg := fmt.Errorf("failed to install package '%s': %s", name, err)
		panic(errorMsg)
	}

	fmt.Println(string(stdout))
}
