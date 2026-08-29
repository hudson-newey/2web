package lexer

import (
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
)

func textLexer(model *Lexer) (V2LexNode, LexFunc) {
	cases := lexDefMap{
		// The comment state has to always come first that it takes precedence over
		// other matches and can omit them as source text.
		"<!--": {token: lexeme.MarkupCommentStart, next: markupCommentLexer},

		"<":  {token: lexeme.LessAngle, next: elementLexer},
		">":  {token: lexeme.GreaterAngle, next: textLexer},
		"\\": {token: lexeme.Escape, next: textLexer},

		"{": {token: lexeme.TextContent, next: textLexer},
		"}": {token: lexeme.TextContent, next: textLexer},

		// "\n": {token: lexeme.NewLine, next: textLexer},
		// "\t": {token: lexeme.Tab, next: textLexer},
	}

	return lexerFactory(cases, textContent)(model)
}
