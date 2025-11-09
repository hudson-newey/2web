package scanners

import (
	"fmt"
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
)

func FirstToken(tokens []*lexer.V2LexNode, matcher lexerTokens.LexToken) (*lexer.V2LexNode, error) {
	for _, candidate := range tokens {
		if candidate.Token == matcher {
			return candidate, nil
		}
	}

	return nil, fmt.Errorf("token not found. Expected: '%s'", matcher.String())
}
