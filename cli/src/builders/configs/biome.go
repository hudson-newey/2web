package configs

import (
	"errors"
	"os"
)

func BiomeConfigLocation() (string, error) {
	overridePath := "biome.json"
	if _, err := os.Stat(overridePath); err == nil {
		return overridePath, nil
	}

	sdkPath := "node_modules/@two-web/cli/templates/biome.json"
	if _, err := os.Stat(sdkPath); err == nil {
		return sdkPath, nil
	}

	return "", errors.New("could not find biome config")
}
