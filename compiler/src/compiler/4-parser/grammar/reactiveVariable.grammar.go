package grammar

import (
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
	"hudson-newey/2web/src/compiler/4-parser/nodes"
)

var scriptVariable = Grammar{
	Def: newDefinition(
		lexeme.DollarSign,
		lexeme.CompiledScriptSource, // variable name
		lexeme.Equals,
		lexeme.CompiledScriptSource, // initial value
		lexeme.Semicolon,
	),
	Constructor: wrapConstructor(nodes.NewreactiveVariableNode),
}
