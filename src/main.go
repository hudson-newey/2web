package main

import (
	"hudson-newey/2web/src/cli"
	"hudson-newey/2web/src/compiler"
	"os"
)

func main() {
	args := cli.ParseArguments()

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
	compilerResult := compiler.CompileFile(inputPath)
	os.WriteFile(outputPath, []byte(compilerResult), 0644)
}
