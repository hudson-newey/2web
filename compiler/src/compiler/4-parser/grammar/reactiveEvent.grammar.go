package grammar

import (
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
	"hudson-newey/2web/src/compiler/4-parser/nodes"
)

var reactiveEvent = Grammar{
	Def: newDefinition(
		lexerTokens.AtSymbol,
		lexerTokens.TextContent,
		lexerTokens.Equals,
		lexerTokens.QuoteDouble,
		lexerTokens.TextContent,
	),
	Constructor: wrapConstructor(nodes.NewReactiveEventNode),
}
