package grammar

import (
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
	"hudson-newey/2web/src/compiler/4-parser/nodes"
)

var inlineScripts = Grammar{
	Def: newDefinition(
		lexeme.LessAngle,
		lexeme.ScriptStartTag,
		lexeme.GreaterAngle,
		lexeme.ScriptSource,
		lexeme.ScriptEndTag,
	),
	Constructor: wrapConstructor(nodes.NewScriptNode),
	ChildDefs:   []Grammar{},
}
