package grammar

import lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"

type definition = []lexerTokens.LexToken

type Grammar struct {
	// A sequence of tokens that define the reactive variable
	Def definition

	// A constructor function to create a node from the tokens
	Ctor func(...lexerTokens.LexToken) any
}
