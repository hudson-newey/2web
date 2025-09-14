package lexer

import (
	"hudson-newey/2web/src/compiler/2-lexer/states"
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
)

// When inside the first starting angle bracket (<) and up until (and including)
// the closing angle bracket (>).
func elementLexer(model *Lexer) (V2LexNode, LexFunc) {
	cases := lexDefMap{
		"!doctype": {token: lexerTokens.Doctype, next: elementLexer},

		// I treat tabs like spaces so that they are treated the same in attributes
		" ":  {token: lexerTokens.Space, next: elementLexer},
		"\t": {token: lexerTokens.Space, next: elementLexer},
		"/":  {token: lexerTokens.Slash, next: elementLexer},
		"'":  {token: lexerTokens.QuoteSingle, next: elementLexer},
		"\"": {token: lexerTokens.QuoteDouble, next: elementLexer},
		"@":  {token: lexerTokens.AtSymbol, next: elementLexer},
		"*":  {token: lexerTokens.Star, next: elementLexer},
		"#":  {token: lexerTokens.Hash, next: elementLexer},
		"=":  {token: lexerTokens.Equals, next: elementLexer},
		"!":  {token: lexerTokens.Exclamation, next: textLexer},
		">":  {token: lexerTokens.GreaterAngle, next: textLexer},
	}

	cases = withStrings(cases, elementLexer)

	return lexerFactory(cases, states.Element)(model)
}
