package cache

import (
	"fmt"
	"hash/adler32"
	"hudson-newey/2web/src/constants"
	"os"
	"path"
	"time"
)

func cacheRecordLocation(
	inputPath string,
	outputPath string,
	inputModified time.Time,
) string {
	directoryPath := cacheLocation()

	// TODO: Use an sqlite database instead of checksums here
	inputNameHash := adler32.Checksum([]byte(inputPath))
	outputNameHash := adler32.Checksum([]byte(outputPath))

	// We use the modified time instead of a file because it is easier to compute
	// and can also be invalidated by using the "touch" command.
	fileName := fmt.Sprintf(
		"%x-%x-%x",
		inputNameHash,
		outputNameHash,
		inputModified.Unix(),
	)

	return path.Join(directoryPath, fileName)
}

func cacheLocation() string {
	currentDir, _ := os.Getwd()
	defaultCachePath := path.Join(currentDir, "/.cache/")

	overrideValue, hasOverride := os.LookupEnv(constants.EnvCacheOverride)
	if hasOverride {
		return overrideValue
	}

	return defaultCachePath
}
