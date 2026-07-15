package lexer

import (
	"hudson-newey/2web/src/compiler/2-lexer/states"
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
)

// Attribute lexing is shared between multiple states.
// e.g. The inline style and script tags can have attribute but require their
// own lexer because any text state they transition into is actually a new lexer
// state where you are lexing script/css text rather than normal html text.
func withAttributes(src lexDefMap) lexDefMap {
	attributeStates := lexDefMap{
		// I treat tabs like spaces so that they are treated the same in attributes
		" ":  {token: lexerTokens.Space, next: elementLexer},
		"\t": {token: lexerTokens.Space, next: elementLexer},
		"/":  {token: lexerTokens.Slash, next: elementLexer},
		"'":  {token: lexerTokens.QuoteSingle, next: elementLexer},
		"\"": {token: lexerTokens.QuoteDouble, next: elementLexer},
		"#":  {token: lexerTokens.Hash, next: elementLexer},
		"=":  {token: lexerTokens.Equals, next: elementLexer},
		"!":  {token: lexerTokens.Exclamation, next: textLexer},
		">":  {token: lexerTokens.GreaterAngle, next: textLexer},

		"*": {token: lexerTokens.Star, next: reactivePropertyLexer},
		"@": {token: lexerTokens.AtSymbol, next: reactiveEventLexer},
	}

	return src.with(attributeStates)
}

func reactivePropertyLexer(model *Lexer) (V2LexNode, LexFunc) {
	cases := lexDefMap{
		"=": {token: lexerTokens.Equals, next: reactivePropertyLexer},
	}
	cases = withStrings(cases, elementLexer)
	return lexerFactory(cases, states.Element)(model)
}

func reactiveEventLexer(model *Lexer) (V2LexNode, LexFunc) {
	cases := lexDefMap{
		"=": {token: lexerTokens.Equals, next: reactiveEventLexer},
	}
	cases = withStrings(cases, elementLexer)
	return lexerFactory(cases, states.Element)(model)
}
