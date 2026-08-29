package lexer

import (
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
)

func inlineCodeTagLexer(model *Lexer) (V2LexNode, LexFunc) {
	cases := lexDefMap{
		">": {token: lexeme.GreaterAngle, next: codeContentLexer},
	}

	cases = withAttributes(cases)
	cases = withStrings(cases, inlineCodeTagLexer)

	return lexerFactory(cases, codeSource)(model)
}

func codeContentLexer(model *Lexer) (V2LexNode, LexFunc) {
	cases := lexDefMap{
		"</code>": {token: lexeme.CodeEndTag, next: textLexer},
	}

	return lexerFactory(cases, codeSource)(model)
}
