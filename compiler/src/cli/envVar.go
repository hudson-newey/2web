package cli

import (
	"fmt"
	"hudson-newey/2web/src/constants"
	"os"
	"path"
	"sync"
)

type envVars struct {
	CacheOverride string
	DebugOverride string
	IsCi          bool
}

// cachedEnvVars memoizes the resolved environment variables.
//
// GetEnvVars is consulted from per file hot paths (the cache location is
// resolved for every cache operation), and resolving it performs an os.Getwd
// syscall. The environment and CLI arguments are constant for the lifetime of
// the process, so the result is computed once.
var envVarsOnce sync.Once
var cachedEnvVars envVars

func GetEnvVars() envVars {
	envVarsOnce.Do(func() {
		cachedEnvVars = resolveEnvVars()
	})

	return cachedEnvVars
}

func resolveEnvVars() envVars {
	cacheOverride, hasOverride := os.LookupEnv(constants.EnvCacheOverride)
	if !hasOverride {
		currentDir, _ := os.Getwd()
		defaultCachePath := path.Join(currentDir, "/.cache/")
		cacheOverride = defaultCachePath
	}

	debugOverride, hasOverride := os.LookupEnv(constants.EnvDebugOverride)
	if !hasOverride {
		outputDirectory := GetArgs().OutputPath

		// An empty output path would panic on the index below. It can't be
		// suffixed with a trailing slash either, so fall back to the working
		// directory.
		if outputDirectory == "" {
			outputDirectory = "./"
		} else if outputDirectory[len(outputDirectory)-1] != '/' {
			outputDirectory += "/"
		}

		// Default to {output_path}/__2web.debug.json so that it's accessible
		// to third party extensions that look at the build output.
		debugOverride = fmt.Sprintf("%s/__2web.debug.json", outputDirectory)
	}

	isCiString, hasOverride := os.LookupEnv(constants.EnvCiOverride)
	if !hasOverride {
		isCiString = "false"
	}

	// Use "true" as the default so that we prefer "false" values over "true"
	// any value except for "true" will result in a non-CI environment since I
	// suspect more non-ci environments than CI.
	// Therefore, targeting non-CI environments is more important.
	isCiBool := isCiString == "true"

	return envVars{
		CacheOverride: cacheOverride,
		DebugOverride: debugOverride,
		IsCi:          isCiBool,
	}
}
