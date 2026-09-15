package lexer

import (
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
)

var scriptCommentLexerState stateLexers

func scriptCommentLexer(model *Lexer) (V2LexNode, LexFunc) {
	matchers := scriptCommentLexerState.get(func() lexDefMap {
		return lexDefMap{
			"*/": {token: lexeme.MarkupCommentStart, next: inlineScriptTagLexer},
		}
	})

	return lexerFactory(matchers, scriptComment)(model)
}
