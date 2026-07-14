package grammar

import (
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
	"hudson-newey/2web/src/compiler/4-parser/nodes"
)

var compiledScripts = Grammar{
	Def: newDefinition(
		lexerTokens.LessAngle,
		lexerTokens.CompiledScriptStartTag,
		lexerTokens.GreaterAngle,
		lexerTokens.NewCaptureUntil(),
		lexerTokens.ScriptEndTag,
	),
	Constructor: wrapConstructor(nodes.NewTwoScriptNode),
	ChildDefs: []Grammar{
		scriptImport,
		scriptVariable,
	},
}
