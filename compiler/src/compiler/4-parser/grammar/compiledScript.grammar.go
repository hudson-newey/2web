package grammar

import (
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
	"hudson-newey/2web/src/compiler/4-parser/nodes"
)

var compiledScripts = Grammar{
	Def: newDefinition(
		lexeme.LessAngle,
		lexeme.CompiledScriptStartTag,
		lexeme.GreaterAngle,
		lexeme.NewCaptureUntil(),
		lexeme.ScriptEndTag,
	),
	Constructor: wrapConstructor(nodes.NewTwoScriptNode),
	ChildDefs: []Grammar{
		scriptImport,
		scriptVariable,
	},
}
