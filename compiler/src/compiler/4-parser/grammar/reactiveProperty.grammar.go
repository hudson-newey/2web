package grammar

import (
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
	"hudson-newey/2web/src/compiler/4-parser/nodes"
)

var reactiveProperty = Grammar{
	Def: newDefinition(
		lexerTokens.Star,
		lexerTokens.TextContent,
		lexerTokens.Equals,
		lexerTokens.QuoteDouble,
		lexerTokens.TextContent,
	),
	Constructor: wrapConstructor(nodes.NewReactivePropertyNode),
}
