package lexer

import (
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
	"hudson-newey/2web/src/compiler/2-lexer/states"
)

func scriptCommentLexer(model *Lexer) (V2LexNode, LexFunc) {
	cases := lexDefMap{
		"*/": {token: lexeme.MarkupCommentStart, next: inlineScriptTagLexer},
	}

	return lexerFactory(cases, states.ScriptComment)(model)
}
