package cli

import (
	"fmt"
	"hudson-newey/2web/src/constants"
	"os"
	"path"
)

type envVars struct {
	CacheOverride string
	DebugOverride string
	IsCi          bool
}

func GetEnvVars() envVars {
	cacheOverride, hasOverride := os.LookupEnv(constants.EnvCacheOverride)
	if !hasOverride {
		currentDir, _ := os.Getwd()
		defaultCachePath := path.Join(currentDir, "/.cache/")
		cacheOverride = defaultCachePath
	}

	debugOverride, hasOverride := os.LookupEnv(constants.EnvDebugOverride)
	if !hasOverride {
		outputDirectory := GetArgs().OutputPath
		if outputDirectory[len(outputDirectory)-1] != '/' {
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
		IsCi: isCiBool,
	}
}
