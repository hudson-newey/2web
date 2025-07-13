package cli

import (
	"flag"
	"hudson-newey/2web/src/models"
)

var parsedArgs models.CliArguments

func ParseArguments() models.CliArguments {
	inputPath := flag.String("i", "index.html", "Input file path")
	outputPath := flag.String("o", "./dist/index.html", "Output file path")
	hasDevTools := flag.Bool("dev-tools", false, "Include dev tools in the build")
	isProd := flag.Bool("production", false, "Optimize code at the cost of readability")
	isSilent := flag.Bool("silent", false, "Do not output log information")
	disableCache := flag.Bool("no-cache", false, "Do not use a build cache")
	fromStdin := flag.Bool("stdin", false, "Read from stdin")
	toStdout := flag.Bool("stdout", false, "Output the build file to stdout instead of to a file location")

	flag.Parse()

	parsedArgs = models.CliArguments{
		InputPath:    inputPath,
		OutputPath:   outputPath,
		HasDevTools:  hasDevTools,
		IsProd:       isProd,
		IsSilent:     isSilent,
		DisableCache: disableCache,
		FromStdin:    fromStdin,
		ToStdout:     toStdout,
	}

	return GetArgs()
}

func GetArgs() models.CliArguments {
	return parsedArgs
}
