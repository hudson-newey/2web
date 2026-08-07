package grammar

import (
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
	"hudson-newey/2web/src/compiler/4-parser/nodes"
)

var inlineStyles = Grammar{
	Def: newDefinition(
		lexeme.LessAngle,
		lexeme.StyleStartTag,
		lexeme.GreaterAngle,
		lexeme.StyleSource,
		lexeme.StyleEndTag,
	),
	Constructor: wrapConstructor(nodes.NewStyleNode),
	ChildDefs:   []Grammar{},
}
