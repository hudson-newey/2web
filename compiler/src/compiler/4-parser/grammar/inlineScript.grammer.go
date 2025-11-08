package grammar

import (
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
	"hudson-newey/2web/src/compiler/4-parser/nodes"
)

var inlineScripts = Grammar{
	Def: definition{
		lexerTokens.LessAngle,
		lexerTokens.ScriptStartTag,
		lexerTokens.GreaterAngle,
		lexerTokens.ScriptSource,
		lexerTokens.ScriptEndTag,
	},
	Constructor: wrapConstructor(nodes.NewScriptNode),
}
