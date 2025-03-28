package cli

import (
	"flag"
	"hudson-newey/2web/src/models"
)

func ParseArguments() models.CliArguments {
	inputPath := flag.String("i", "index.html", "Input file path")
	outputPath := flag.String("o", "./dist/index.html", "Output file path")

	flag.Parse()

	return models.CliArguments{
		InputPath:  inputPath,
		OutputPath: outputPath,
	}
}
