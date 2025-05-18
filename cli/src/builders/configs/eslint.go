package configs

import (
	"errors"
	"os"
)

func EslintConfigLocation() string {
	overridePath := "eslintrc.js"
	if _, err := os.Stat(overridePath); errors.Is(err, os.ErrNotExist) {
		return "node_modules/@two-web/cli/sdk/eslintrc.js"
	}

	return overridePath
}
