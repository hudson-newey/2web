package lexer

import (
	"hudson-newey/2web/src/compiler/2-lexer/states"
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
)

// Because string lexers are dynamically based on the starting context and quote
// type, we use a factory for two reasons:
//  1. To create a new lexer function that has the correct returning state
//  2. To have the correct exiting token based on the quote type
func createStringLexer(returningState LexFunc, exitToken LexMatcher) LexFunc {
	// TODO: If the user escapes the quote type, we want to ignore it and continue
	// processing the string.
	cases := lexDefMap{
		exitToken: {token: lexerTokens.LexToken(exitToken), next: returningState},
	}

	return lexerFactory(cases, states.String)
}

// Since strings are common to most states, this helper function adds all of the
// string related transitions to an existing lexDefMap.
func withStrings(src lexDefMap, returningState LexFunc) lexDefMap {
	stringStates := lexDefMap{
		`"`: {token: lexerTokens.QuoteDouble, next: createStringLexer(returningState, "\"")},
		`'`: {token: lexerTokens.QuoteSingle, next: createStringLexer(returningState, "'")},
		"`": {token: lexerTokens.Backtick, next: createStringLexer(returningState, "`")},
	}

	return src.with(stringStates)
}
