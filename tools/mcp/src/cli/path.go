package cli

import (
	"os"
	"os/exec"
	"path/filepath"
)

// Because 2web projects sometimes have a local bin/ directory in the projects
// root directory, we need to conditionally return the local bin/2web path if
// it exists, otherwise return the global 2web CLI path.
// This function can throw an error if the 2web CLI is not installed at the
// project level or globally.
func twoWebCliPath() (string, error) {
	localPath, err := filepath.Abs("./bin/2web")
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(localPath); err == nil {
		return localPath, nil
	}

	globalPath, err := exec.LookPath("2web")
	if globalPath == "" || err != nil {
		return "", os.ErrNotExist
	}

	return globalPath, nil
}
