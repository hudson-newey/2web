package builder

import (
	"hudson-newey/2web/src/builder/cache"
	"hudson-newey/2web/src/cli"
	"hudson-newey/2web/src/content/page/runtimeOptimizer"
	"hudson-newey/2web/src/optimizer"
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
	if !cacheDisabled && !production {
		hasCached := cache.IsCached(inputPath, outputPath)
		if hasCached {
			cli.PrintBuildLog("\t- " + inputPath + " \033[36m(cached)\033[0m")
			return
		}
	}

	compiledPage, success := BuildToPage(inputPath, true)

	// We perform runtime optimizations here so that they are only applied once.
	// If we instead applied them inside of the BuildToPage function, then the
	// optimizations would be applied every time a component is added.
	// Most optimizations can be applied at the very end of the page build, so
	// this is the best place to do it.
	if !(args.NoRuntimeOptimizations) {
		runtimeOptimizer.InjectRuntimeOptimizations(&compiledPage)
	}

	if args.WithFormatting {
		compiledPage.Format()
	}

	// We always optimize last so that even the injected content is optimized.
	if args.IsProd {
		if args.WithFormatting {
			cli.PrintWarning("Ignoring '--format' because '--production' was specified")
		}

		optimizer.OptimizePage(&compiledPage)
	}

	if !success && production && !args.IgnoreErrors {
		// Compiler errors should not be ignored in production builds, otherwise, we
		// start shipping compiler errors to end users, which does not look good.
		cli.HardError("Build failed for " + inputPath)
		return
	}

	compiledPage.WriteHtml(outputPath)

	if success {
		if cacheDisabled {
			cli.PrintBuildLog("\t- " + inputPath)
		} else {
			cli.PrintBuildLog("\t- " + inputPath + " \033[33m(MODIFIED)\033[0m")
		}

		if !cacheDisabled {
			cache.CacheAsset(inputPath, outputPath)
		}
	} else {
		cli.PrintBuildLog("\t- " + inputPath + " \033[31m(ERROR)\033[0m")
	}
}
