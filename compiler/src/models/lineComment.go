package models

import "hudson-newey/2web/src/compiler/lexer"

type LineComment struct {
	Node *lexer.LexNode[lexer.LineCommentNode]
}
