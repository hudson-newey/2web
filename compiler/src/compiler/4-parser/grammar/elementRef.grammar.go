package grammar

import (
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
	"hudson-newey/2web/src/compiler/4-parser/nodes"
)

var elementRef = Grammar{
	// #myElement
	Def: newDefinition(
		lexeme.Hash,
		lexeme.TextContent,
	),
	Constructor: wrapConstructor(nodes.NewElementRefNode),
}
