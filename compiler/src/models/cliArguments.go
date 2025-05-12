package models

type CliArguments struct {
	InputPath  *string
	OutputPath *string
	IsDev      *bool
	IsProd     *bool
	IsSilent   *bool
	ToStdout   *bool
}
