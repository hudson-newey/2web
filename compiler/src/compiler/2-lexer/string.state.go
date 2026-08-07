package lexer

import (
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
	"hudson-newey/2web/src/compiler/2-lexer/states"
)

// Because string lexers are dynamically based on the starting context and quote
// type, we use a factory for two reasons:
//  1. To create a new lexer function that has the correct returning state
//  2. To have the correct exiting token based on the quote type
func createStringLexer(returningState LexFunc, exitToken LexMatcher) LexFunc {
	// TODO: If the user escapes the quote type, we want to ignore it and continue
	// processing the string.
	cases := lexDefMap{
		exitToken: {token: lexeme.Lexeme(exitToken), next: returningState},
	}

	return lexerFactory(cases, states.String)
}

// Since strings are common to most states, this helper function adds all of the
// string related transitions to an existing lexDefMap.
func withStrings(src lexDefMap, returningState LexFunc) lexDefMap {
	stringStates := lexDefMap{
		`"`: {token: lexeme.QuoteDouble, next: createStringLexer(returningState, "\"")},
		`'`: {token: lexeme.QuoteSingle, next: createStringLexer(returningState, "'")},
		"`": {token: lexeme.Backtick, next: createStringLexer(returningState, "`")},
	}

	return src.with(stringStates)
}
