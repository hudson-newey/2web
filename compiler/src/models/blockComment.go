package models

import lexer "hudson-newey/2web/src/compiler/2-lexer"

type BlockComment struct {
	Node *lexer.LexNode[lexer.BlockCommentNode]
}
