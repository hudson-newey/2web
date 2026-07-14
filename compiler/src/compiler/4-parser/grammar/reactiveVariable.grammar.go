package grammar

import (
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
	"hudson-newey/2web/src/compiler/4-parser/nodes"
)

var scriptVariable = Grammar{
	Def: newDefinition(
		lexerTokens.DollarSign,
		lexerTokens.CompiledScriptSource, // variable name
		lexerTokens.Equals,
		lexerTokens.CompiledScriptSource, // initial value
		lexerTokens.Semicolon,
	),
	Constructor: wrapConstructor(nodes.NewScriptReactiveVariableNode),
}
