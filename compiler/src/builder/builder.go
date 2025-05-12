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
	"log"
	"os"
	"path/filepath"
)

func Build(args models.CliArguments) {
	if *args.IsDev && *args.IsProd {
		cli.PrintWarning("'--dev-tools' is being used with '--production'")
	}

	log.Println(*args.InputPath)

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

			compileAndWriteFile(*args.InputPath+"/"+file.Name(), *args.OutputPath+"/"+file.Name(), &args)
		}
	} else {
		compileAndWriteFile(*args.InputPath, *args.OutputPath, &args)
	}
}

func compileAndWriteFile(
	inputPath string,
	outputPath string,
	args *models.CliArguments,
) {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		data = []byte{}
		documentErrors.AddError(models.Error{
			FilePath: inputPath,
			Message:  fmt.Sprintf("Failed to read file: %s\n%s", inputPath, err.Error()),
		})
	}

	log.Println("\t-", inputPath)

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

	writeOutput(finalResult, outputPath, args)
}

func writeOutput(content string, outputPath string, args *models.CliArguments) {
	if *args.ToStdout {
		fmt.Println(content)
	} else {
		os.MkdirAll(filepath.Dir(outputPath), os.ModePerm)
		os.WriteFile(outputPath, []byte(content), 0644)
	}
}
