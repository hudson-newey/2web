package lexer

import "hudson-newey/2web/src/compiler/2-lexer/lexeme"

var controlFlowLexerState stateLexers

func controlFlowLexer(model *Lexer) (V2LexNode, LexFunc) {
	matchers := controlFlowLexerState.get(func() lexDefMap {
		return lexDefMap{
			"if": {token: lexeme.ControlFlowIfKeyword, next: controlFlowLexer},
			"(":  {token: lexeme.BracketOpen, next: controlFlowLexer},
			")":  {token: lexeme.BracketClosed, next: controlFlowLexer},

			// exit condition is a closing curly brace
			"{": {token: lexeme.CurlyOpen, next: controlFlowLexer},
			"}": {token: lexeme.CurlyClosed, next: textLexer},
		}
	})

	return lexerFactory(matchers, controlFlow)(model)
}
