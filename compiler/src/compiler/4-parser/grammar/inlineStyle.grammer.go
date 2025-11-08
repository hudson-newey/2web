package grammar

import (
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
	"hudson-newey/2web/src/compiler/4-parser/nodes"
)

var inlineStyles = Grammar{
	Def: definition{
		lexerTokens.LessAngle,
		lexerTokens.StyleStartTag,
		lexerTokens.GreaterAngle,
		lexerTokens.StyleSource,
		lexerTokens.StyleEndTag,
	},
	Constructor: wrapConstructor(nodes.NewStyleNode),
}
