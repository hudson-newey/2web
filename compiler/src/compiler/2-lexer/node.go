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

type V2LexNode struct {
	Pos     Position
	Token   lexerTokens.LexToken
	State   states.LexState
	Content string
}

func (model *V2LexNode) PrintDebug() string {
	// We replace all of the new lines and tabs with their escape character
	// representations so that the output is easier to read.
	quotedContent := strconv.Quote(model.Content)

	return fmt.Sprintf("%d:%d\t%s\t%s\t%s\n", model.Pos.Row, model.Pos.Col, model.State, model.Token, quotedContent)
}
