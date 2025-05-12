package builder

import (
	"fmt"
	"hudson-newey/2web/src/cli"
	"hudson-newey/2web/src/compiler"
	"hudson-newey/2web/src/compiler/controlFlow"
	"hudson-newey/2web/src/document/devtools"
	"hudson-newey/2web/src/document/documentErrors"
	"hudson-newey/2web/src/models"
	"hudson-newey/2web/src/optimizer"
	"hudson-newey/2web/src/ssg"
	"io"
	"os"
	"path/filepath"
)

func Build() {
	args := cli.GetArgs()

	if *args.IsDev && *args.IsProd {
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

			compileAndWriteFile(*args.InputPath+"/"+file.Name(), *args.OutputPath+"/"+file.Name())
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

	controlFlowResult := controlFlow.ProcessControlFlow(inputPath, string(data))

	ssgSource := controlFlowResult
	stable := false
	for {
		ssgSource, stable = ssg.ProcessStaticSite(inputPath, ssgSource)

		if stable {
			break
		}
	}

	compilerResult := compiler.Compile(inputPath, ssgSource)

	injectedErrorResult := documentErrors.InjectErrors(compilerResult)

	finalResult := injectedErrorResult
	if *args.IsDev {
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
