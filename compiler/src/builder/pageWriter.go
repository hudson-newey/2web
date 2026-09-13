package builder

import (
	"hudson-newey/2web/src/builder/cache"
	"hudson-newey/2web/src/cli"
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/content/document/documentErrors"
	"hudson-newey/2web/src/models"
	"hudson-newey/2web/src/site"
)

// Writes an entry point to the output path.
func compileAndWritePage(inputPath string, outputPath string) {
	args := cli.GetArgs()
	cacheDisabled := args.DisableCache
	production := args.IsProd

	site.BeforeEach(outputPath)

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
	//
	// The cache key is a fingerprint of the page's source and all of its
	// compile time dependencies (layouts, sidecars, and imported components and
	// ESM modules), so a change to any of them invalidates the page.
	//
	// The key is computed once here and reused for the cache registration so
	// that the recorded cache entry always matches the inputs that this page
	// was actually compiled against.
	cacheKey := ""
	if !cacheDisabled && !production {
		computedKey, err := cache.CacheKey(inputPath, outputPath)
		if err != nil {
			cli.PrintWarning("could not compute build cache key for " + inputPath + ": " + err.Error())
			computedKey = ""
		}

		cacheKey = computedKey

		if cacheKey != "" && cache.IsCached(inputPath, outputPath, cacheKey) {
			cli.PrintBuildLog("\t- " + inputPath + " \033[36m(cached)\033[0m")
			return
		}
	}

	compiledPage, success := BuildToPage(inputPath, true)

	if !success && production && !args.IgnoreErrors {
		// Compiler errors should not be ignored in production builds, otherwise,
		// we start shipping compiler errors to end users, which does not look
		// good.
		//
		// The failure is recorded as a compiler error (which is rendered into
		// the page and reported in the terminal) and the build continues so that
		// the remaining pages still compile. The build is marked as failed
		// through the document error list.
		buildError := models.NewError(
			"build failed (production builds do not ship compiler errors)",
			inputPath,
			lexer.StartingPosition,
		)

		documentErrors.AddErrors(&buildError)
		cli.PrintBuildLog("\t- " + inputPath + " \033[31m(FAILED)\033[0m")
		return
	}

	compiledPage.WriteHtml(outputPath)

	if success {
		if cacheDisabled {
			cli.PrintBuildLog("\t- " + inputPath)
		} else {
			cli.PrintBuildLog("\t- " + inputPath + " \033[33m(MODIFIED)\033[0m")
			// We don't add the asset to the build cache here because the page's
			// files are written asynchronously. queueCacheAsset defers the cache
			// registration until the file writes have been confirmed, which
			// prevents the cache from marking assets as up-to-date when their
			// output file was never written.
			queueCacheAsset(inputPath, outputPath, cacheKey)
		}
	} else {
		cli.PrintBuildLog("\t- " + inputPath + " \033[31m(ERROR)\033[0m")
	}
}
