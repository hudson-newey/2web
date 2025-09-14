package parser

import (
	"fmt"
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/content/document/documentErrors"
	"hudson-newey/2web/src/models"
)

type parserError = int

const (
	syntaxError parserError = iota
)

func printError(lexer lexer.Lexer, node lexer.V2LexNode, err parserError) {
	message := fmt.Sprintf("Error (%d:%d): %s", node.Pos.Row, node.Pos.Col, errorMessage(err))
	errorModel := models.Error{
		Position: node.Pos,
		FilePath: lexer.Input.FilePath,
		Message:  message,
	}

	documentErrors.AddErrors(errorModel)
}

func errorMessage(err parserError) string {
	switch err {
	case syntaxError:
		return "Syntax error"
	default:
		return "Unknown error"
	}
}
