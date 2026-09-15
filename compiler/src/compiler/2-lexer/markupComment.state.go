package lexer

import (
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
)

var markupCommentLexerState stateLexers

func markupCommentLexer(model *Lexer) (V2LexNode, LexFunc) {
	matchers := markupCommentLexerState.get(func() lexDefMap {
		return lexDefMap{
			"-->": {token: lexeme.MarkupCommentEnd, next: textLexer},
		}
	})

	return lexerFactory(matchers, markupComment)(model)
}
