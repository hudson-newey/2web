package lexer

import "hudson-newey/2web/src/compiler/2-lexer/lexeme"

var controlFlowLexerState stateLexers

func controlFlowLexer(model *Lexer) (V2LexNode, LexFunc) {
	matchers := controlFlowLexerState.get(func() lexDefMap {
		return lexDefMap{
			"if":  {token: lexeme.ControlFlowIfKeyword, next: controlFlowLexer},
			"for": {token: lexeme.ControlFlowForKeyword, next: controlFlowLexer},
			"(":   {token: lexeme.BracketOpen, next: controlFlowLexer},
			")":   {token: lexeme.BracketClosed, next: controlFlowLexer},

			// A doubled closing curly brace is matched as one token (before the
			// single closing brace; matchers are ordered by length) so that a
			// control flow body can contain text outputs like '{{ $item }}'
			// without the body capture ending at the text output's own brace.
			"}}": {token: lexeme.DoubleCurlyClosed, next: controlFlowLexer},

			// exit condition is a closing curly brace
			"{": {token: lexeme.CurlyOpen, next: controlFlowLexer},
			"}": {token: lexeme.CurlyClosed, next: textLexer},
		}
	})

	return lexerFactory(matchers, controlFlow)(model)
}
