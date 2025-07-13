package builder

import (
	"hudson-newey/2web/src/cli"
	"hudson-newey/2web/src/content/document/documentErrors"
	"hudson-newey/2web/src/models"
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
		documentErrors.AddErrors(models.Error{
			FilePath: *args.InputPath,
			Message:  err.Error(),
		})

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

		for _, file := range indexedPages {
			compileAndWriteFile(file, outputFileName(*args.InputPath, *args.OutputPath, file))
		}
	} else {
		compileAndWriteFile(*args.InputPath, outputFileName(*args.InputPath, *args.OutputPath, *args.InputPath))
	}

	// We only print document errors once the entire project has been compiled so
	// that all of the errors are located in the same section, and are the most
	// recent output when compilation finishes.
	documentErrors.PrintDocumentErrors()

	return documentErrors.IsErrorFree()
}
