package configs

import (
	"errors"
	"os"
)

// Gets the "web test runner" config location
func WtrConfigLocation() (string, error) {
	overridePaths := []string{
		"web-test-runner.config.js",
		"web-test-runner.config.cjs",
		"web-test-runner.config.mjs",
	}

	for _, override := range overridePaths {
		if _, err := os.Stat(override); err == nil {
			return override, nil
		}
	}

	sdkPath := "node_modules/@two-web/cli/templates/web-test-runner.config.mjs"
	if _, err := os.Stat(sdkPath); err == nil {
		return sdkPath, nil
	}

	return "", errors.New("could not find web test runner config")
}
