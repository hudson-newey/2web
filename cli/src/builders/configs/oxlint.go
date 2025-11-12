package configs

import (
	"errors"
	"os"
)

func OxLintConfigLocation() (string, error) {
	overridePath := ".oxlintrc.json"
	if _, err := os.Stat(overridePath); err == nil {
		return overridePath, nil
	}

	sdkPath := "node_modules/@two-web/cli/templates/.oxlintrc.json"
	if _, err := os.Stat(sdkPath); err == nil {
		return sdkPath, nil
	}

	return "", errors.New("could not find oxlint config")
}
