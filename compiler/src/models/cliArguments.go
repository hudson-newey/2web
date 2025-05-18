package models

type CliArguments struct {
	InputPath   *string
	OutputPath  *string
	HasDevTools *bool
	IsProd      *bool
	IsSilent    *bool
	FromStdin   *bool
	ToStdout    *bool
}
