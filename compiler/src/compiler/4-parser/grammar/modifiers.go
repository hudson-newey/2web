package grammar

import "hudson-newey/2web/src/compiler/2-lexer/lexeme"

// optional wraps a lexeme so that a grammar definition matches zero or one
// occurrence of it.
//
// This is commonly used to make whitespace between two tokens optional, since
// the lexer emits whitespace as a TextContent token:
//
//	@if ($isOpen) { ... }
//	 ^^              ^^
//	 \______________/\__ two TextContent tokens that should not be required
func optional(def lexeme.Lexeme) lexeme.Lexeme {
	return lexeme.NewOptional(def)
}

// or creates a grammar token that matches exactly one occurrence of any one of
// the given lexemes.
func or(tokens ...lexeme.Lexeme) lexeme.Lexeme {
	return lexeme.NewOr(tokens...)
}

var anyQuote = or(lexeme.QuoteDouble, lexeme.QuoteSingle)
