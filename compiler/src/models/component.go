package models

import "hudson-newey/2web/src/compiler/lexer"

type Component struct {
	// The selector that can be used in the template to reference this component
	// e.g. "Footer" for <Footer />
	DomSelector string
	ImportPath  string
	Node        *lexer.LexNode[lexer.ImportNode]
}
