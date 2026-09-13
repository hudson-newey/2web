package lexer

import (
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
)

var inlineCodeTagLexerState stateLexers

func inlineCodeTagLexer(model *Lexer) (V2LexNode, LexFunc) {
	matchers := inlineCodeTagLexerState.get(func() lexDefMap {
		cases := lexDefMap{
			">": {token: lexeme.GreaterAngle, next: codeContentLexer},
		}

		cases = withAttributes(cases, inlineCodeTagLexer, "inline-code-tag")
		cases = withStrings(cases, inlineCodeTagLexer, "inline-code-tag")
		return cases
	})

	return lexerFactory(matchers, tagAttributes)(model)
}

var codeContentLexerState stateLexers

func codeContentLexer(model *Lexer) (V2LexNode, LexFunc) {
	matchers := codeContentLexerState.get(func() lexDefMap {
		return lexDefMap{
			"</code>": {token: lexeme.CodeEndTag, next: textLexer},
		}
	})

	return lexerFactory(matchers, codeSource)(model)
}
