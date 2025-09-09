package lexer

import (
	"hudson-newey/2web/src/compiler/2-lexer/states"
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
)

func styleLexer(model *Lexer) (V2LexNode, LexFunc) {
	cases := lexDefMap{
		"</style>": {token: lexerTokens.NewLine, next: textLexer},
	}

	cases = withStrings(cases, styleLexer)

	return lexerFactory(cases, states.StyleSource)(model)
}
