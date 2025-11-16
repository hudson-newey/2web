package grammar

import (
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
	"hudson-newey/2web/src/compiler/4-parser/nodes"
)

var scriptVariable = grammar{
	Def: newDefinition(
		lexerTokens.DollarSign,
		lexerTokens.TextContent, // variable name
		lexerTokens.Equals,
		lexerTokens.TextContent, // initial value
		lexerTokens.Semicolon,
	),
	Constructor: wrapConstructor(nodes.NewScriptReactiveVariableNode),
	ChildDefs:   []grammar{},
}
