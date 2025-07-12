package validator

import (
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/models"
)

func IsValid(structure lexer.LexerRepresentation) (bool, []models.Error) {
	return true, []models.Error{}
}
