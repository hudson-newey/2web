package grammar

import "hudson-newey/2web/src/compiler/2-lexer/lexeme"

var anyQuote = or(lexeme.QuoteDouble, lexeme.QuoteSingle)

func optional(def lexeme.Lexeme) lexeme.Lexeme {
	return def
}

func or(tokens ...lexeme.Lexeme) lexeme.Lexeme {
	return tokens[0]
}
