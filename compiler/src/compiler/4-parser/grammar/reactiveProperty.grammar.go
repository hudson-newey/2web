package grammar

import (
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
	"hudson-newey/2web/src/compiler/4-parser/nodes"
)

var reactiveProperty = Grammar{
	Def: newDefinition(
		lexeme.Star,
		lexeme.TextContent,
		lexeme.Equals,
		lexeme.QuoteDouble,
		lexeme.TextContent,
	),
	Constructor: wrapConstructor(nodes.NewReactivePropertyNode),
}
