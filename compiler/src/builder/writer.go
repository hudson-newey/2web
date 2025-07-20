package builder

import (
	"hudson-newey/2web/src/builder/cache"
	"hudson-newey/2web/src/cli"
)

func compileAndWriteFile(inputPath string, outputPath string) {
	args := cli.GetArgs()
	cacheDisabled := *args.DisableCache
	production := *args.IsProd

	// We keep a record of the last modified time of all recent input files, so
	// that we can skip re-compiling source files that have not changed.
	//
	// We have to pass in the output path so that we can validate that the output
	// file wasn't deleted by the user or an external process.
	//
	// Additionally, production build speeds should be concerned with
	// correctness instead of build speed.
	// We therefore, skip the build cache if we are building for production, so
	// that cache errors are (almost) impossible in production builds.
	//
	// Additionally, production builds have some additional processing that dev
	// builds do not.
	// Therefore, our method of checking modified times does not work if the
	// build output is different.
	if !cacheDisabled && !production {
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
