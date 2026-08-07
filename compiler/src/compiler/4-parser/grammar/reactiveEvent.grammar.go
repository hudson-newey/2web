package grammar

import (
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
	"hudson-newey/2web/src/compiler/4-parser/nodes"
)

var reactiveEvent = Grammar{
	Def: newDefinition(
		lexeme.AtSymbol,
		lexeme.TextContent,
		lexeme.Equals,
		lexeme.QuoteDouble,
		lexeme.TextContent,
	),
	Constructor: wrapConstructor(nodes.NewReactiveEventNode),
}
