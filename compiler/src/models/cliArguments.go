package models

type CliArguments struct {
	InputPath              *string
	OutputPath             *string
	HasDevTools            *bool
	NoRuntimeOptimizations *bool
	IsProd                 *bool
	IsSilent               *bool
	DisableCache           *bool
	FromStdin              *bool
	ToStdout               *bool

	// I only expect end users to use the "Verbose" cli argument.
	// However, I have still exposed lower-level verbosity levels so that debug
	// information can also be generated from end-users.
	Verbose      *bool
	VerboseLexer *bool
	VerboseAst   *bool
}
