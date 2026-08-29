package lexer

import (
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
)

// Attribute lexing is shared between multiple states.
// e.g. The inline style and script tags can have attribute but require their
// own lexer because any text state they transition into is actually a new lexer
// state where you are lexing script/css text rather than normal html text.
func withAttributes(src lexDefMap) lexDefMap {
	attributeStates := lexDefMap{
		// I treat tabs like spaces so that they are treated the same in attributes
		" ":  {token: lexeme.Space, next: elementLexer},
		"\t": {token: lexeme.Space, next: elementLexer},
		"/":  {token: lexeme.Slash, next: elementLexer},
		"'":  {token: lexeme.QuoteSingle, next: elementLexer},
		"\"": {token: lexeme.QuoteDouble, next: elementLexer},
		"#":  {token: lexeme.Hash, next: elementLexer},
		"=":  {token: lexeme.Equals, next: elementLexer},
		"!":  {token: lexeme.Exclamation, next: textLexer},
		">":  {token: lexeme.GreaterAngle, next: textLexer},

		"*": {token: lexeme.Star, next: reactivePropertyLexer},
		"@": {token: lexeme.AtSymbol, next: reactiveEventLexer},
	}

	return src.with(attributeStates)
}

func reactivePropertyLexer(model *Lexer) (V2LexNode, LexFunc) {
	cases := lexDefMap{
		"=": {token: lexeme.Equals, next: reactivePropertyLexer},
	}
	cases = withStrings(cases, elementLexer)
	return lexerFactory(cases, element)(model)
}

func reactiveEventLexer(model *Lexer) (V2LexNode, LexFunc) {
	cases := lexDefMap{
		"=": {token: lexeme.Equals, next: reactiveEventLexer},
	}
	cases = withStrings(cases, elementLexer)
	return lexerFactory(cases, element)(model)
}
