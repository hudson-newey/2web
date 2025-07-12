package builder

import (
	"fmt"
	"hudson-newey/2web/src/cli"
	preprocessor "hudson-newey/2web/src/compiler/1-preprocessor"
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	validator "hudson-newey/2web/src/compiler/3-validator"
	templating "hudson-newey/2web/src/compiler/5-templating"
	"hudson-newey/2web/src/content/document/devtools"
	"hudson-newey/2web/src/content/document/documentErrors"
	"hudson-newey/2web/src/content/markdown"
	"hudson-newey/2web/src/content/page"
	"hudson-newey/2web/src/models"
	"hudson-newey/2web/src/optimizer"
	"io"
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

func buildToPage(inputPath string) page.Page {
	args := cli.GetArgs()

	rawData, err := getInputContent(inputPath)
	if err != nil {
		rawData = []byte{}
		documentErrors.AddErrors(models.Error{
			FilePath: inputPath,
			Message:  fmt.Sprintf("Failed to read file: %s\n%s", inputPath, err.Error()),
			Position: models.Position{
				Line: 1,
			},
		})
	}

	cli.PrintBuildLog("\t- " + inputPath)

	ssgResult := preprocessor.ProcessStaticSite(inputPath, string(rawData))

	pageModel := templating.BuildPage(ssgResult)

	lexInstance := lexer.NewLexer(pageModel.Html.Reader())
	lexRepresentation := lexInstance.Execute()

	isValid, compilerErrors := validator.IsValid(lexRepresentation)
	if !isValid {
		documentErrors.AddErrors(compilerErrors...)
	}

	compiledPage := templating.Compile(inputPath, pageModel)

	compiledPage.Html.Content = documentErrors.InjectErrors(compiledPage.Html.Content)
	documentErrors.ResetPageErrors()

	if *args.HasDevTools {
		compiledPage.Html.Content = devtools.InjectDevTools(compiledPage.Html.Content)
	}

	if *args.IsProd {
		compiledPage = optimizer.OptimizePage(compiledPage)
	}

	return compiledPage
}

func compileAndWriteFile(inputPath string, outputPath string) {
	compiledPage := buildToPage(inputPath)
	compiledPage.Write(outputPath)
}

func getInputContent(inputPath string) ([]byte, error) {
	if !*cli.GetArgs().FromStdin {
		return os.ReadFile(inputPath)
	}

	if !*cli.GetArgs().IsSilent {
		fmt.Println("Prompting STDIN for file:", inputPath)
	}

	return io.ReadAll(os.Stdin)
}
