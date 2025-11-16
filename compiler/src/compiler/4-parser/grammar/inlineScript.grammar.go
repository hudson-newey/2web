package grammar

import (
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
	"hudson-newey/2web/src/compiler/4-parser/nodes"
)

var inlineScripts = grammar{
	Def: newDefinition(
		lexerTokens.LessAngle,
		lexerTokens.GreaterAngle,
		lexerTokens.ScriptSource,
		lexerTokens.ScriptEndTag,
	),
	Constructor: wrapConstructor(nodes.NewScriptNode),
	ChildDefs: []grammar{
		scriptImport,
		scriptVariable,
	},
}
