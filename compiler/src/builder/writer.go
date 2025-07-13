package builder

import (
	"hudson-newey/2web/src/builder/cache"
	"hudson-newey/2web/src/cli"
)

func compileAndWriteFile(inputPath string, outputPath string) {
	cacheDisabled := *cli.GetArgs().DisableCache

	// We keep a record of the last modified time of all recent input files, so
	// that we can skip re-compiling source files that have not changed.
	//
	// We have to pass in the output path so that we can validate that the output
	// file wasn't deleted by the user or an external process.
	if !cacheDisabled {
		hasCached := cache.IsCached(inputPath, outputPath)
		if hasCached {
			cli.PrintBuildLog("\t- " + inputPath + " \033[36m(cached)\033[0m")
			return
		}
	}

	compiledPage, success := buildToPage(inputPath)
	compiledPage.Write(outputPath)

	if success {
		cli.PrintBuildLog("\t- " + inputPath)

		if !cacheDisabled {
			cache.CacheAsset(inputPath, outputPath)
		}
	} else {
		cli.PrintBuildLog("\t- " + inputPath + " \033[31m(error)\033[0m")
	}
}
