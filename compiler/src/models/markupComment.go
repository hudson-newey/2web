package models

import lexer "hudson-newey/2web/src/compiler/2-lexer"

// Markup comments can be used within html, svg, MathML, markdown and other
// browser-based markup languages.
type MarkupComment struct {
	Node *lexer.LexNode[lexer.LineCommentNode]
}
