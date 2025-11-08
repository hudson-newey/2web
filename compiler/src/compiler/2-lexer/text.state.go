package lexer

import (
	"hudson-newey/2web/src/compiler/2-lexer/states"
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
)

func textLexer(model *Lexer) (V2LexNode, LexFunc) {
	cases := lexDefMap{
		// The comment state has to always come first that it takes precedence over
		// other matches and can omit them as source text.
		"<!--": {token: lexerTokens.MarkupCommentStart, next: markupCommentLexer},

		"<":  {token: lexerTokens.LessAngle, next: elementLexer},
		">":  {token: lexerTokens.GreaterAngle, next: textLexer},
		"\\": {token: lexerTokens.Escape, next: textLexer},

		"{": {token: lexerTokens.TextContent, next: textLexer},
		"}": {token: lexerTokens.TextContent, next: textLexer},

		// "\n": {token: lexerTokens.NewLine, next: textLexer},
		// "\t": {token: lexerTokens.Tab, next: textLexer},
	}

	return lexerFactory(cases, states.TextContent)(model)
}
