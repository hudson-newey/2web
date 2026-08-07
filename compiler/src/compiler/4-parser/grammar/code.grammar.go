package grammar

import (
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
	"hudson-newey/2web/src/compiler/4-parser/nodes"
)

var codeBlock = Grammar{
	Def: newDefinition(
		lexeme.LessAngle,
		lexeme.CodeStartTag,
		lexeme.GreaterAngle,
		lexeme.CodeSource,
		lexeme.CodeEndTag,
	),
	Constructor: wrapConstructor(nodes.NewCodeNode),
}
