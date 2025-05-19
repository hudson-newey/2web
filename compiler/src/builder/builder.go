package builder

import (
	"fmt"
	"hudson-newey/2web/src/cli"
	"hudson-newey/2web/src/compiler/ssg"
	"hudson-newey/2web/src/compiler/templating"
	"hudson-newey/2web/src/compiler/templating/controlFlow"
	"hudson-newey/2web/src/content/document/devtools"
	"hudson-newey/2web/src/content/document/documentErrors"
	"hudson-newey/2web/src/content/html"
	"hudson-newey/2web/src/content/markdown"
	"hudson-newey/2web/src/models"
	"hudson-newey/2web/src/optimizer"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Build() {
	args := cli.GetArgs()

	if *args.HasDevTools && *args.IsProd {
		cli.PrintWarning("'--dev-tools' is being used with '--production'")
	}

	cli.PrintBuildLog(*args.InputPath)

	inputPath, err := os.Stat(*args.InputPath)
	if err != nil {
		panic(err)
	}

	if inputPath.IsDir() {
		// find all direct children of the input directory
		files, err := os.ReadDir(*args.InputPath)
		if err != nil {
			panic(err)
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			// If we are compiling a markdown file, we want to replace the .md suffix
			// with .html
			// This is because we compile markdown to html files.
			adjustedFileName := file.Name()
			if markdown.IsMarkdownFile(file.Name()) {
				adjustedFileName = strings.TrimSuffix(adjustedFileName, ".md") + ".html"
			}

			compileAndWriteFile(*args.InputPath+"/"+file.Name(), *args.OutputPath+"/"+adjustedFileName)
		}
	} else {
		compileAndWriteFile(*args.InputPath, *args.OutputPath)
	}
}

func compileAndWriteFile(inputPath string, outputPath string) {
	args := cli.GetArgs()

	data, err := getInputContent(inputPath)
	if err != nil {
		data = []byte{}
		documentErrors.AddError(models.Error{
			FilePath: inputPath,
			Message:  fmt.Sprintf("Failed to read file: %s\n%s", inputPath, err.Error()),
		})
	}

	cli.PrintBuildLog("\t- " + inputPath)

	fullDocumentContent := ""
	if html.IsHtmlFile(inputPath) {
		// 2Web supports partial content, meaning that pages don't need and doctype,
		// html, head, meta, or body tags.
		// The user can just start writing the pages content, and the compiler can
		// figure out what should be in the body vs head.
		fullDocumentContent = html.ExpandPartial(string(data))
	} else if markdown.IsMarkdownFile(inputPath) {
		markdownFile := markdown.MarkdownFile{
			Content: string(data),
		}
		fullDocumentContent = markdownFile.ToHtml().Content

		// Markdown files are typically compiled as html partials. Developers
		// typically (and shouldn't) declare a doctype, header, etc...
		// therefore, we also expand html partials once the html document has been
		// created.
		fullDocumentContent = html.ExpandPartial(fullDocumentContent)
	} else {
		fullDocumentContent = string(data)
	}

	controlFlowResult := controlFlow.ProcessControlFlow(inputPath, fullDocumentContent)

	ssgSource := controlFlowResult
	stable := false
	for {
		ssgSource, stable = ssg.ProcessStaticSite(inputPath, ssgSource)

		if stable {
			break
		}
	}

	compilerResult := templating.Compile(inputPath, ssgSource)

	injectedErrorResult := documentErrors.InjectErrors(compilerResult)

	finalResult := injectedErrorResult
	if *args.HasDevTools {
		finalResult = devtools.InjectDevTools(injectedErrorResult)
	}

	if *args.IsProd {
		finalResult = optimizer.OptimizeContent(finalResult)
	}

	writeOutput(finalResult, outputPath)
}

func writeOutput(content string, outputPath string) {
	if *cli.GetArgs().ToStdout {
		fmt.Println(content)
	} else {
		os.MkdirAll(filepath.Dir(outputPath), os.ModePerm)
		os.WriteFile(outputPath, []byte(content), 0644)
	}
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
