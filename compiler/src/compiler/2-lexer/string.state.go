package lexer

import (
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
	"sync"
)

// stringLexerCache caches the compiled matcher set of every string lexer
// variant (per calling state and quote type) so that quoted attribute values
// don't rebuild their matchers on every quote token.
var stringLexerCache sync.Map // string -> LexFunc

// Because string lexers are dynamically based on the starting context and quote
// type, we use a factory for two reasons:
//  1. To create a new lexer function that has the correct returning state
//  2. To have the correct exiting token based on the quote type
func createStringLexer(returningState LexFunc, exitToken LexMatcher, cacheKey string) LexFunc {
	cacheLookup := cacheKey + "\x00" + string(exitToken)

	if cached, ok := stringLexerCache.Load(cacheLookup); ok {
		return cached.(LexFunc)
	}

	var once sync.Once
	var matchers *compiledMatchers

	lexerFunc := func(model *Lexer) (V2LexNode, LexFunc) {
		once.Do(func() {
			matchers = compileMatchers(lexDefMap{
				exitToken: {token: lexeme.Lexeme(exitToken), next: returningState},
			})
		})

		return lexerFactory(matchers, sourceString)(model)
	}

	stringLexerCache.Store(cacheLookup, lexerFunc)
	return lexerFunc
}

// Since strings are common to most states, this helper function adds all of the
// string related transitions to an existing lexDefMap.
//
// The cacheKey identifies the calling state so that the compiled string
// matchers (which capture the calling state's return transition) can be cached
// per caller.
func withStrings(src lexDefMap, returningState LexFunc, cacheKey string) lexDefMap {
	stringStates := lexDefMap{
		`"`: {token: lexeme.QuoteDouble, next: createStringLexer(returningState, "\"", cacheKey)},
		`'`: {token: lexeme.QuoteSingle, next: createStringLexer(returningState, "'", cacheKey)},
		"`": {token: lexeme.Backtick, next: createStringLexer(returningState, "`", cacheKey)},
	}

	return src.with(stringStates)
}
