package lexer

import lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"

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
		"@":  {token: lexerTokens.AtSymbol, next: elementLexer},
		"*":  {token: lexerTokens.Star, next: elementLexer},
		"#":  {token: lexerTokens.Hash, next: elementLexer},
		"=":  {token: lexerTokens.Equals, next: elementLexer},
		"!":  {token: lexerTokens.Exclamation, next: textLexer},
		">":  {token: lexerTokens.GreaterAngle, next: textLexer},
	}

	return src.with(attributeStates)
}
