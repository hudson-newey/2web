package builder

import (
	"fmt"
	"hudson-newey/2web/src/compiler"
	"hudson-newey/2web/src/document/devtools"
	"hudson-newey/2web/src/document/documentErrors"
	"hudson-newey/2web/src/models"
	"hudson-newey/2web/src/ssg"
	"log"
	"os"
)

func Build(args models.CliArguments) {
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

			compileAndWriteFile(*args.InputPath+"/"+file.Name(), *args.OutputPath+"/"+file.Name(), *args.IsDev)
		}
	} else {
		compileAndWriteFile(*args.InputPath, *args.OutputPath, *args.IsDev)
	}
}

func compileAndWriteFile(inputPath string, outputPath string, isDev bool) {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		data = []byte{}
		documentErrors.AddError(models.Error{
			FilePath: inputPath,
			Message:  fmt.Sprintf("Failed to read file: %s\n%s", inputPath, err.Error()),
		})
	}

	log.Println("\t-", inputPath)

	ssgSource := string(data)
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
	if isDev {
		finalResult = devtools.InjectDevTools(injectedErrorResult)
	}

	os.WriteFile(outputPath, []byte(finalResult), 0644)
}
