package lexer

import (
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
)

func markupCommentLexer(model *Lexer) (V2LexNode, LexFunc) {
	cases := lexDefMap{
		"-->": {token: lexeme.MarkupCommentEnd, next: textLexer},
	}

	return lexerFactory(cases, markupComment)(model)
}
