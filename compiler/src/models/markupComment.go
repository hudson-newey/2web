package models

import "hudson-newey/2web/src/compiler/lexer"

// Markup comments can be used within html, svg, MathML, markdown and other
// browser-based markup languages.
type MarkupComment struct {
	Node *lexer.LexNode[lexer.LineCommentNode]
}
