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
	noRuntimeOptimizations := flag.Bool("skip-rasterization-optimizations", false, "Slightly reduces bundle size, but makes page (e.g. scrolling) less responsive")
	isProd := flag.Bool("production", false, "Optimize code at the cost of readability")
	isSilent := flag.Bool("silent", false, "Do not output log information")
	disableCache := flag.Bool("no-cache", false, "Do not use a build cache")
	fromStdin := flag.Bool("stdin", false, "Read from stdin")
	toStdout := flag.Bool("stdout", false, "Output the build file to stdout instead of to a file location")
	withFormatting := flag.Bool("format", false, "Formats output assets for readability")
	ignoreErrors := flag.Bool("ignore-errors", false, "Ignores errors in production builds. This allows you to ship compiler errors.")

	verbose := flag.Bool("verbose", false, "Print extra debug information to the console")
	verboseLexer := flag.Bool("verbose-lexer", false, "Logs the lexer output to the console")
	verboseAst := flag.Bool("verbose-ast", false, "Logs the ast output to the console")

	flag.Parse()

	parsedArgs = models.CliArguments{
		InputPath:              inputPath,
		OutputPath:             outputPath,
		HasDevTools:            hasDevTools,
		NoRuntimeOptimizations: noRuntimeOptimizations,
		IsProd:                 isProd,
		IsSilent:               isSilent,
		DisableCache:           disableCache,
		FromStdin:              fromStdin,
		ToStdout:               toStdout,
		WithFormatting:         withFormatting,
		IgnoreErrors:           ignoreErrors,

		Verbose:      verbose,
		VerboseLexer: verboseLexer,
		VerboseAst:   verboseAst,
	}

	return GetArgs()
}

func GetArgs() models.CliArguments {
	return parsedArgs
}
