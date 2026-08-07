package grammar

import "hudson-newey/2web/src/compiler/2-lexer/lexeme"

type definition []lexeme.Lexeme

func newDefinition(tokens ...lexeme.Lexeme) definition {
	return definition(tokens)
}
