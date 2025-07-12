package models

import "hudson-newey/2web/src/compiler/2-lexer"

type LineComment struct {
	Node *lexer.LexNode[lexer.LineCommentNode]
}
