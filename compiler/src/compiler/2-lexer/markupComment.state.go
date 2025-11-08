package lexer

import (
	"hudson-newey/2web/src/compiler/2-lexer/states"
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
)

func markupCommentLexer(model *Lexer) (V2LexNode, LexFunc) {
	cases := lexDefMap{
		"-->": {token: lexerTokens.MarkupCommentEnd, next: textLexer},
	}

	return lexerFactory(cases, states.MarkupComment)(model)
}
