package lexer

import (
	"hudson-newey/2web/src/compiler/2-lexer/states"
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
)

func inlineCodeTagLexer(model *Lexer) (V2LexNode, LexFunc) {
	cases := lexDefMap{
		">": {token: lexerTokens.GreaterAngle, next: codeContentLexer},
	}

	cases = withAttributes(cases)
	cases = withStrings(cases, inlineCodeTagLexer)

	return lexerFactory(cases, states.CodeSource)(model)
}

func codeContentLexer(model *Lexer) (V2LexNode, LexFunc) {
	cases := lexDefMap{
		"</code>": {token: lexerTokens.CodeEndTag, next: textLexer},
	}

	return lexerFactory(cases, states.CodeSource)(model)
}
