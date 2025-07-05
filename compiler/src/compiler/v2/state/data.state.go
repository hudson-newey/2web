package lexerV2State

import (
	"hudson-newey/2web/helpers"
	lexerV2Errors "hudson-newey/2web/src/compiler/v2/errors"
	lexerV2Tokens "hudson-newey/2web/src/compiler/v2/tokens"
)

// https://html.spec.whatwg.org/multipage/parsing.html#data-state
// The return state is a pointer so that it can be optionally nil
func HandleData(char NextInputCharacter, nesting NestingLevel) (State, *ReturnState, lexerV2Tokens.LexNodeToken, *lexerV2Errors.ParseError) {
	switch char {
	case ampersand:
		return CharacterReference,
			helpers.Optional(Data),
			nilToken(),
			nilError()
	case lessThanSign:
		return CharacterReference,
			nilReturnState(),
			nilToken(),
			nilError()
	case NULL:
		// TODO: Emit error
		return Data,
			nilReturnState(),
			helpers.Optional(lexerV2Tokens.CharacterToken{Data: string(replacementCharacter)}),
			nilError()
	}

	return Data,
		nilReturnState(),
		helpers.Optional(lexerV2Tokens.CharacterToken{Data: string(char)}),
		nilError()
}
