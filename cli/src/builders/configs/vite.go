package configs

import (
	"errors"
	"os"
)

func ViteConfigLocation() string {
	// If the override path exists, the user has explicitly created a
	// vite.config.ts file in the project root, and we should use this config
	// instead of using the SDK Vite config.
	overridePath := "vite.config.ts"
	if _, err := os.Stat(overridePath); errors.Is(err, os.ErrNotExist) {
		return "node_modules/@two-web/cli/sdk/vite.config.ts"
	}

	return overridePath
}
