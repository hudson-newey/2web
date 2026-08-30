package lexer

import "hudson-newey/2web/src/compiler/2-lexer/lexeme"

func controlFlowLexer(model *Lexer) (V2LexNode, LexFunc) {
	cases := lexDefMap{
		"if": {token: lexeme.ControlFlowIfKeyword, next: controlFlowLexer},
		"(":  {token: lexeme.BracketOpen, next: controlFlowLexer},
		")":  {token: lexeme.BracketClosed, next: controlFlowLexer},

		// exit condition is a closing curly brace
		"{": {token: lexeme.CurlyOpen, next: controlFlowLexer},
		"}": {token: lexeme.CurlyClosed, next: textLexer},
	}

	return lexerFactory(cases, controlFlow)(model)
}
