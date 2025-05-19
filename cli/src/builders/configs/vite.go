package configs

import (
	"errors"
	"os"
)

func ViteConfigLocation() (string, error) {
	// If the override path exists, the user has explicitly created a
	// vite.config.ts file in the project root, and we should use this config
	// instead of using the SDK Vite config.
	overridePath := "vite.config.ts"
	if _, err := os.Stat(overridePath); err == nil {
		return overridePath, nil
	}

	sdkPath := "node_modules/@two-web/sdk/vite.config.ts"
	if _, err := os.Stat(sdkPath); err == nil {
		return sdkPath, nil
	}

	// If an empty path is returned, we can determine that a vite config does not
	// exist in the current project, either in the current directory or in the
	// sdk path.
	// We therefore use the default config provided by Vite.
	return "", errors.New("could not find vite config")
}
