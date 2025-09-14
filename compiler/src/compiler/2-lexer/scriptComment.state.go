package lexer

import (
	"hudson-newey/2web/src/compiler/2-lexer/states"
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
)

func scriptCommentLexer(model *Lexer) (V2LexNode, LexFunc) {
	cases := lexDefMap{
		"*/": {token: lexerTokens.MarkupCommentStart, next: scriptLexer},
	}

	return lexerFactory(cases, states.ScriptComment)(model)
}
