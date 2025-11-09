package cache

import (
	"hudson-newey/2web/src/constants"
	"os"
	"path"
)

func cacheLocation() string {
	currentDir, _ := os.Getwd()
	defaultCachePath := path.Join(currentDir, "/.cache/")

	overrideValue, hasOverride := os.LookupEnv(constants.EnvCacheOverride)
	if hasOverride {
		return overrideValue
	}

	return defaultCachePath
}
