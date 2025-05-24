package configs

import (
	"errors"
	"os"
)

func EslintConfigLocation() (string, error) {
	overridePath := "eslintrc.js"
	if _, err := os.Stat(overridePath); err == nil {
		return overridePath, nil
	}

	sdkPath := "node_modules/@two-web/sdk/eslintrc.js"
	if _, err := os.Stat(sdkPath); err == nil {
		return sdkPath, nil
	}

	return "", errors.New("could not find eslint config")
}
