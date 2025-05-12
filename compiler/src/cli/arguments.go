package cli

import (
	"flag"
	"hudson-newey/2web/src/models"
)

func ParseArguments() models.CliArguments {
	inputPath := flag.String("i", "index.html", "Input file path")
	outputPath := flag.String("o", "./dist/index.html", "Output file path")
	isDev := flag.Bool("dev-tools", false, "Include dev tools in the build")
	isProd := flag.Bool("production", false, "Optimize code at the cost of readability")
	isSilent := flag.Bool("silent", false, "Do not output log information")
	toStdout := flag.Bool("stdout", false, "Output the build file to stdout instead of to a file location")

	flag.Parse()

	return models.CliArguments{
		InputPath:  inputPath,
		OutputPath: outputPath,
		IsDev:      isDev,
		IsProd:     isProd,
		IsSilent:   isSilent,
		ToStdout:   toStdout,
	}
}
