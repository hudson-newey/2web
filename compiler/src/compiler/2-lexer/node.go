package lexer

import (
	"fmt"
	"hudson-newey/2web/src/compiler/2-lexer/states"
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
	"strconv"
)

type LexNodeType[T voidNode] any

type Position struct {
	Row int
	Col int
}

func NewV2LexNode() V2LexNode {
	return V2LexNode{}
}

type V2LexNode struct {
	Pos     Position
	Token   lexerTokens.LexToken
	State   states.LexState
	Content string
}

// The 0th position in a file (row 0, column 0).
// This is NOT the same as the position of the first character in a file
// (which is typically row 1, column 1).
var StartingPosition = Position{Row: 0, Col: 0}

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
