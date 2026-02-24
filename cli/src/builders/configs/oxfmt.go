package configs

import (
	"errors"
	"os"
)

func OxFmtConfigLocation() (string, error) {
	overridePath := ".oxfmtrc.json"
	if _, err := os.Stat(overridePath); err == nil {
		return overridePath, nil
	}

	sdkPath := "./node_modules/@two-web/cli/templates/.oxfmtrc.json"
	if _, err := os.Stat(sdkPath); err == nil {
		return sdkPath, nil
	}

	return "", errors.New("could not find oxfmt config")
}
