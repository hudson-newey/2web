package lexer

import (
	"fmt"
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
	"hudson-newey/2web/src/models"
	"strconv"
)

type LexNodeType[T voidNode] any

// Position is a location in a source file.
//
// It is an alias of the models package position so that the lexer and the
// error models can share the type without an import cycle.
type Position = models.Position

func NewV2LexNode() V2LexNode {
	return V2LexNode{}
}

type V2LexNode struct {
	Pos     Position
	Token   lexeme.Lexeme
	State   lexState
	Content string
}

// The 0th position in a file (row 0, column 0).
// This is NOT the same as the position of the first character in a file
// (which is typically row 1, column 1).
var StartingPosition = models.StartingPosition

func (model *V2LexNode) PrintDebug() string {
	// We replace all of the new lines and tabs with their escape character
	// representations so that the output is easier to read.
	quotedContent := strconv.Quote(model.Content)

	return fmt.Sprintf(
		"%d:%d\t%s\t%s\t%s\n",
		model.Pos.Row,
		model.Pos.Col,
		model.State.String(),
		model.Token.String(),
		quotedContent,
	)
}
