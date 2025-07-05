package lexerV2State

import (
	"hudson-newey/2web/helpers"
	lexerV2Errors "hudson-newey/2web/src/compiler/v2/errors"
	lexerV2Tokens "hudson-newey/2web/src/compiler/v2/tokens"
)

func nilReturnState() *ReturnState {
	return helpers.Optional[ReturnState]()
}

func nilToken() *lexerV2Tokens.LexNodeToken {
	return helpers.Optional[lexerV2Tokens.LexNodeToken]()
}

func nilError() *lexerV2Errors.ParseError {
	return helpers.Optional[lexerV2Errors.ParseError]()
}
