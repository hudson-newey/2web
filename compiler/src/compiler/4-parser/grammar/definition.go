package grammar

import (
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
)

type definition []lexerTokens.LexToken

func newDefinition(tokens ...lexerTokens.LexToken) definition {
	return definition(tokens)
}
