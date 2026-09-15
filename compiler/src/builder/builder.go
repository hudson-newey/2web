package builder

import (
	"fmt"
	"hudson-newey/2web/src/builder/cache"
	"hudson-newey/2web/src/cli"
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/content/document/documentErrors"
	"hudson-newey/2web/src/debugger"
	"hudson-newey/2web/src/filesystem"
	"hudson-newey/2web/src/models"
	"hudson-newey/2web/src/site"
	"os"
	"time"

	"github.com/hudson-newey/2web/_shared/logger"
)

func Build() bool {
	startTime := time.Now()

	// Every cache touch of this build is recorded with the build's start time
	// (see cache.BeginBuild), so the cache writes never need to query the
	// clock themselves.
	cache.BeginBuild(startTime)

	args := cli.GetArgs()

	if args.HasDevTools && args.IsProd {
		cli.PrintWarning("'--dev-tools' is being used with '--production'")
	}

	cli.PrintBuildLog(args.InputPath)

	inputPath, err := os.Stat(args.InputPath)
	if err != nil {
		pathError := models.NewError(err.Error(), args.InputPath, lexer.Position{})
		documentErrors.AddErrors(&pathError)

		documentErrors.PrintDocumentErrors()
		return false
	}

	// If the output path already exists, delete the output path so that there
	// are no stale files.
	if _, err := os.Stat(args.OutputPath); err == nil {
		// TODO: This has been disabled because it doesn't work with Vite HMR
		// os.RemoveAll(args.OutputPath)
	}

	// Print out the "starting compilation" message after we have confirmed that
	// the build command is constructed correctly.
	// Because builds may progressively invalidate some output file integrity,
	// the build message serves as a "output files can not be trusted message to
	// the developer which is incorrect if we said this message while still
	// parsing arguments.
	if !cli.GetArgs().IsSilent {
		logger.Println("Begin 2webc compilation...")
	}

	filesystem.InitFileWriter()
	startBuildPool()

	if inputPath.IsDir() {
		// recursively find all children of the input directory
		indexedPages := indexPages(args.InputPath)

		// The debug file writer carries reactive graphs over from previous
		// builds for pages that this build served from cache. Recording the
		// input files that belong to this build lets it prune graphs that were
		// carried over for pages that aren't part of the site anymore (or that
		// belong to a different input directory).
		debugger.SetCurrentBuildPages(indexedPages)

		for _, filePath := range indexedPages {
			AddCompilationStep(filePath)
		}

		buildPool.Wait()

		// TODO: These AfterAll hooks should use the output asset models as the
		// input so if any files implicitly create any site-wide assets, the
		// AfterAll hook does not overwrite the files.
		//
		// We only perform a site-wide AfterAll hook if we are compiling a page of
		// results so that you don't end up with a bunch of unnecessary site-wide
		// files when compiling a single document.
		// Note that compiling a single document is primarily used by external tools
		// such as bundler plugins, linters, etc...
		site.AfterAll()
	} else {
		debugger.SetCurrentBuildPages([]string{args.InputPath})
		compileAndWritePage(args.InputPath, outputFileName(args.InputPath, args.OutputPath, args.InputPath))
	}

	// Page compilation hands file writes off to an asynchronous file writer
	// thread, so reaching this point only means that all pages have been
	// compiled - not that their output has been written to disk.
	//
	// Server routes (.server.ts files) are compiled into the server output
	// directory and mounted on the generated express server. The route
	// manifest is written after every route has been compiled.
	FlushServerRoutes()

	// We must drain the file write queue before reporting the build as
	// finished, otherwise the compiler could exit with writes still sitting in
	// the queue (silently dropping them) or with the writer thread still
	// mid-write (leaving a partially written file on disk).
	filesystem.WaitFileWriter()

	// Now that all file writes have been flushed to disk, we know exactly which
	// outputs were written successfully and can safely record them in the
	// build cache. Entries that this build served from cache are marked as
	// touched with it as well.
	flushCacheEntries()
	cache.FlushTouches()

	// Reclaim space in the build cache when it has grown too large. This only
	// does work once the cache database is actually over its size limit, and
	// the build is finished at this point, so the (rare) vacuum cost isn't
	// part of the compile time the user is waiting on.
	cache.Vacuum()

	// Printing out document errors are not included in the compile time since
	// the app is fully usable at this point.
	//
	// If the user is repeatedly compiling (e.g. through an automated build),
	// then they are probably not watching the compiler output.
	// This means that if the compiler for some reason stops compiling your app
	// we want to somehow notify the developer that the app hasn't re-compiled
	// when they are not watching the terminal.
	// To accomplish this, we add the last compile finish time to the output in
	// the hopes that when they check the compiler output, they'll see the last
	// compiler time was not recent.
	// Because this use case typically only needs ~10 minute accuracy, we don't
	// need to log out the day/month/year format since it would not be very
	// useful.
	//
	// The only use case I can think day/month/year information being useful is
	// for CI/CD environments where maybe there's 1 deploy a day/month.
	// Therefore, if we are running in CI/CD environments, we log out the full
	// day/month/year format so that you can easily see if the last build was
	// not from when you expected it to be.
	// Note: This information is quite redundant because most CI/CD environments
	// still log out build date/time information, but maybe someday this comes
	// in use lol.
	compileTime := time.Since(startTime)
	buildFinishTime := time.Now()

	buildFinishMessage := buildFinishTime.Format(time.TimeOnly)
	if cli.GetEnvVars().IsCi {
		// Include timezone in CI log output since the CI might not be in the
		// same timezone as the user.
		buildFinishMessage = buildFinishTime.Format(time.RFC3339)
	}

	// We only print document errors once the entire project has been compiled so
	// that all of the errors are located in the same section, and are the most
	// recent output when compilation finishes.
	documentErrors.PrintDocumentErrors()

	// We only print out compile times once the entire app as been built.
	if !cli.GetArgs().IsSilent {
		finishMessage := fmt.Sprintf("\nBuild at: %s - Compile time: %s", buildFinishMessage, compileTime)
		logger.Println(finishMessage)
	}

	return documentErrors.IsErrorFree()
}
