package builder

import (
	"hudson-newey/2web/src/cli"
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/content/document/documentErrors"
	"hudson-newey/2web/src/models"
	"hudson-newey/2web/src/parallel"
	"hudson-newey/2web/src/routing"
	"hudson-newey/2web/src/site"
	"os"
)

func Build() bool {
	args := cli.GetArgs()

	if *args.HasDevTools && *args.IsProd {
		cli.PrintWarning("'--dev-tools' is being used with '--production'")
	}

	cli.PrintBuildLog(*args.InputPath)

	inputPath, err := os.Stat(*args.InputPath)
	if err != nil {
		pathError := models.NewError(err.Error(), *args.InputPath, lexer.Position{})
		documentErrors.AddErrors(&pathError)

		documentErrors.PrintDocumentErrors()
		return false
	}

	// If the output path already exists, delete the output path so that there
	// are no stale files.
	if _, err := os.Stat(*args.OutputPath); err != os.ErrNotExist {
		// TODO: This has been disabled because it doesn't work with Vite HMR
		// os.RemoveAll(*args.OutputPath)
	}

	if inputPath.IsDir() {
		// recursively find all children of the input directory
		indexedPages := indexPages(*args.InputPath)

		parallel.ForEach(indexedPages, func(filePath string) {
			if routing.IsLayoutFile(filePath) {
				return
			}

			compileAndWriteFile(
				filePath,
				outputFileName(*args.InputPath, *args.OutputPath, filePath),
			)
		})

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
		compileAndWriteFile(*args.InputPath, outputFileName(*args.InputPath, *args.OutputPath, *args.InputPath))
	}

	// We only print document errors once the entire project has been compiled so
	// that all of the errors are located in the same section, and are the most
	// recent output when compilation finishes.
	documentErrors.PrintDocumentErrors()

	return documentErrors.IsErrorFree()
}
