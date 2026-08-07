package scanners

import (
	"fmt"
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
)

func FirstToken(
	tokens []*lexer.V2LexNode,
	matcher lexeme.Lexeme,
) (*lexer.V2LexNode, error) {
	return NthToken(tokens, matcher, 1)
}

// NthToken finds the nth token of a certain type in a list of tokens.
// N is 1-indexed meaning that if you pass in n=1, it will return the first
// occurrence.
func NthToken(
	tokens []*lexer.V2LexNode,
	matcher lexeme.Lexeme,
	n int,
) (*lexer.V2LexNode, error) {
	count := 0
	for _, candidate := range tokens {
		if candidate.Token == matcher {
			count++

			if count == n {
				return candidate, nil
			}
		}
	}

	err := fmt.Errorf(
		"nth token not found. Expected: \"%s\" to have \"%d\" occurrences",
		matcher.String(),
		n,
	)

	return nil, err
}
