package scanners

import (
	"fmt"
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
	"strings"
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

// CapturedContent returns the concatenated content of every token between the
// first `start` token and the first `end` token that follows it.
//
// This is used to extract the content of capture blocks (see
// lexeme.NewCaptureUntil) from a grammar's look-ahead buffer. Because the
// captured tokens are concatenated, the returned content preserves the source
// formatting of everything inside the two delimiters.
func CapturedContent(
	tokens []*lexer.V2LexNode,
	start lexeme.Lexeme,
	end lexeme.Lexeme,
) (string, error) {
	insideCapture := false
	content := strings.Builder{}

	for _, candidate := range tokens {
		if !insideCapture {
			if candidate.Token == start {
				insideCapture = true
			}

			continue
		}

		if candidate.Token == end {
			return content.String(), nil
		}

		content.WriteString(candidate.Content)
	}

	return "", fmt.Errorf(
		"captured content not found. Expected a \"%s\" token followed by a \"%s\" token",
		start.String(),
		end.String(),
	)
}
