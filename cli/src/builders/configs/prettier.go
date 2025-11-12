package configs

import (
	"errors"
	"os"
)

func PrettierConfigLocation() (string, error) {
	overridePath := ".prettierrc"
	if _, err := os.Stat(overridePath); err == nil {
		return overridePath, nil
	}

	sdkPath := "node_modules/@two-web/cli/templates/.prettierrc"
	if _, err := os.Stat(sdkPath); err == nil {
		return sdkPath, nil
	}

	return "", errors.New("could not find prettier config")
}
