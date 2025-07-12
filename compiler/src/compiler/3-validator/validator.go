package validator

import (
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/models"
)

func IsValid(structure []lexer.V2LexNode) (bool, []models.Error) {
	// for _, x := range structure {
	// 	fmt.Printf("%d:%d, %d, %s\n", x.Pos.Row, x.Pos.Col, x.Token, x.Content)
	// }

	return true, []models.Error{}
}
