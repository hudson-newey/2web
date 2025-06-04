package models

import "hudson-newey/2web/src/compiler/lexer"

type BlockComment struct {
	Node *lexer.LexNode[lexer.BlockCommentNode]
}
