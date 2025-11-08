package lexer

import (
	"hudson-newey/2web/src/compiler/2-lexer/states"
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
)

// The lexer for when the <style> tag has been opened and before the first >
// meaning that we are technically still in an element tag, but when we
// transition to a text content state, we need to start lexing style content.
func inlineStyleTagLexer(model *Lexer) (V2LexNode, LexFunc) {
	// We need these = and > conditions because inline style tags can have element
	// attributes.
	cases := lexDefMap{
		">": {token: lexerTokens.GreaterAngle, next: styleContentLexer},
	}

	cases = withAttributes(cases)
	cases = withStrings(cases, inlineStyleTagLexer)

	return lexerFactory(cases, states.StyleSource)(model)
}

func styleContentLexer(model *Lexer) (V2LexNode, LexFunc) {
	cases := lexDefMap{
		"</style>": {token: lexerTokens.StyleEndTag, next: textLexer},
	}

	return lexerFactory(cases, states.StyleSource)(model)
}
