package grammar

import (
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
	"hudson-newey/2web/src/compiler/4-parser/nodes"
)

var textOutput = Grammar{
	// {{ $name }}
	Def: newDefinition(
		lexeme.CurlyOpen,
		lexeme.CurlyOpen,
		lexeme.TextContent,
		lexeme.CurlyClosed,
		lexeme.CurlyClosed,
	),
	Constructor: wrapConstructor(nodes.NewTextOutputNode),
}
