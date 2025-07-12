package builder

import (
	"hudson-newey/2web/src/cli"
	"hudson-newey/2web/src/content/document/documentErrors"
	"hudson-newey/2web/src/content/markdown"
	"os"
	"strings"
)

func Build() bool {
	args := cli.GetArgs()

	if *args.HasDevTools && *args.IsProd {
		cli.PrintWarning("'--dev-tools' is being used with '--production'")
	}

	cli.PrintBuildLog(*args.InputPath)

	inputPath, err := os.Stat(*args.InputPath)
	if err != nil {
		panic(err)
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
			// If we are compiling a markdown file, we want to replace the .md suffix
			// with .html
			// This is because we compile markdown to html files.
			adjustedFileName := file
			if markdown.IsMarkdownFile(file) {
				adjustedFileName = strings.TrimSuffix(adjustedFileName, ".md") + ".html"
			}

			adjustedFileName = adjustedFileName[len(*args.InputPath):]

			compileAndWriteFile(file, *args.OutputPath+"/"+adjustedFileName)
		}
	} else {
		compileAndWriteFile(*args.InputPath, *args.OutputPath)
	}

	// We only print document errors once the entire project has been compiled so
	// that all of the errors are located in the same section, and are the most
	// recent output when compilation finishes.
	documentErrors.PrintDocumentErrors()

	return documentErrors.IsErrorFree()
}
