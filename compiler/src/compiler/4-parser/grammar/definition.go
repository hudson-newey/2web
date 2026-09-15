package grammar

import (
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
)

type definition []lexeme.Lexeme

func newDefinition(tokens ...lexeme.Lexeme) definition {
	return definition(tokens)
}

// minimumTokenCount returns the smallest number of input tokens that a
// definition can possibly match against.
//
// Optional tokens can be missing from the input, so they are excluded from the
// count. This is used both as a fast "not enough tokens" bail-out and as the
// parser's look-ahead bound.
func (model *definition) minimumTokenCount() int {
	count := 0
	for _, token := range *model {
		if lexeme.IsOptionalToken(token) {
			continue
		}

		count++
	}

	return count
}
