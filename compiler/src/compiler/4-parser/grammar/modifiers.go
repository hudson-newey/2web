package grammar

import lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"

var anyQuote = or(lexerTokens.QuoteDouble, lexerTokens.QuoteSingle)

func optional(def lexerTokens.LexToken) lexerTokens.LexToken {
	return def
}

func or(tokens ...lexerTokens.LexToken) lexerTokens.LexToken {
	return tokens[0]
}
