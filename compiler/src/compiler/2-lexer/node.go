package lexer

import (
	"fmt"
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
)

type LexNodeType[T voidNode] any

type Position struct {
	Row int
	Col int
}

type V2LexNode struct {
	Pos     Position
	Token   lexerTokens.LexToken
	Content string
}

func (model *V2LexNode) PrintDebug() {
	fmt.Printf("(%d:%d) (%d) %s\n", model.Pos.Row, model.Pos.Col, model.Token, model.Content)
}
