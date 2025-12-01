package models

import (
	"fmt"
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"os"
	"strings"
)

func NewError(message, filePath string, position lexer.Position) Error {
	return Error{
		Message:  message,
		FilePath: filePath,
		Position: position,
	}
}

type Error struct {
	Message      string
	FilePath     string
	Position     lexer.Position
	CreationTime string
}

func (model *Error) PrintError() {
	// Because all error messages are indented by one tab width, I replace all of
	// I add a tab character after every new line.
	formattedErrorMessage := strings.ReplaceAll(model.Message, "\n", "\n\t")

	formattedLineNumber := ""
	if model.Position.Row != 0 {
		formattedLineNumber = fmt.Sprintf(
			" (ln: %d, col: %d)",
			model.Position.Row,
			model.Position.Col,
		)
	}

	// Add a double new line at the start so that the error message is
	// emphasized in the compiler logs.
	fmt.Printf("\n\033[31m[Error] \033[36m%s\033[33m%s\033[0m\n", model.FilePath, formattedLineNumber)
	fmt.Printf("\t%s\n", formattedErrorMessage)
}

func (model *Error) AbsoluteFilePath() string {
	// This prefixes the file path with the current working directory.
	// TODO: This is not a permanent solution. We probably want some sort of path
	// struct that handles this for us.
	currentDir, err := os.Getwd()
	if err == nil {
		return fmt.Sprintf("%s/%s", currentDir, model.FilePath)
	}

	return model.FilePath
}
