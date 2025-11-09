package grammar

import (
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
	"hudson-newey/2web/src/compiler/4-parser/nodes"
)

var codeBlock = grammar{
	Def: definition{
		lexerTokens.LessAngle,
		lexerTokens.CodeStartTag,
		lexerTokens.GreaterAngle,
		lexerTokens.CodeSource,
		lexerTokens.CodeEndTag,
	},
	Constructor: wrapConstructor(nodes.NewCodeNode),
}
