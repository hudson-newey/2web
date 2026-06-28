package models

type CliArguments struct {
	InputPath              string
	OutputPath             string
	HasDevTools            bool
	NoRuntimeOptimizations bool
	IsProd                 bool
	IsSilent               bool
	DisableCache           bool
	FromStdin              bool
	ToStdout               bool
	WithFormatting         bool
	IgnoreErrors           bool
	Serial                 bool
	DryRun                 bool
	Verbose                bool

	// Developer command line flags
	VerboseLexer bool
	VerboseAst   bool
	Listen       bool
}
