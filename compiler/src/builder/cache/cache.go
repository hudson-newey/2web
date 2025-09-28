package cache

import (
	"hudson-newey/2web/src/utils"
	"os"
)

// TODO: Refactor this a lot better because it is still possible to break this.
func IsCached(inputPath string, outputPath string) bool {
	_, err := os.Stat(outputPath)
	if os.IsNotExist(err) {
		return false
	}

	inputFileInfo, _ := os.Stat(inputPath)

	recordLocation := cacheRecordLocation(
		inputPath,
		outputPath,
		inputFileInfo.ModTime(),
	)

	_, err = os.Stat(recordLocation)
	return err == nil
}

func CacheAsset(inputPath string, outputPath string) {
	inputFileInfo, err := os.Stat(inputPath)
	if err != nil {
		panic(err)
	}

	recordLocation := cacheRecordLocation(
		inputPath,
		outputPath,
		inputFileInfo.ModTime(),
	)

	utils.CreateFile(recordLocation)
}
