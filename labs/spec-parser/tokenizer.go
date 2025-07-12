package lexerV2

import (
	lexerV2Errors "hudson-newey/2web/src/compiler/v2/errors"
	lexerV2State "hudson-newey/2web/src/compiler/v2/state"
	lexerV2Tokens "hudson-newey/2web/src/compiler/v2/tokens"
)

// https://html.spec.whatwg.org/multipage/parsing.html#overview-of-the-parsing-model
// https://html.spec.whatwg.org/multipage/parsing.html#tokenization
var scriptNestingLevel = 0
var parserPauseFlag = false
var currentState lexerV2State.State = lexerV2State.Data

// TODO: This should probably become a generic so that if can store different
// types.
// At the moment, this only supports storing return state for character
// reference state.
//
// A buffer that can be used to store temporary state for lexing
// E.g. to store character reference state
// https://html.spec.whatwg.org/multipage/parsing.html#character-reference-state
var tempBuffer lexerV2State.State

func Tokenize(inputData string) []lexerV2Tokens.LexNodeToken {
	handlers := map[lexerV2State.State]lexerV2State.HandlerReturn{
		lexerV2State.Data: lexerV2State.HandleData,
	}

	tokens := []lexerV2Tokens.LexNodeToken{}

	for _, char := range inputData {
		var newToken lexerV2Tokens.LexNodeToken
		var newError *lexerV2Errors.ParseError
		var newTempBuffer *lexerV2State.State

		currentState, newTempBuffer, newToken, newError = handlers[currentState](char, scriptNestingLevel)

		if newTempBuffer != nil {
			tempBuffer = *newTempBuffer
		}

		if newToken != nil {
			tokens = append(tokens, newToken)
		}

		if newError != nil {
			panic("failed to parse html")
		}
	}

	return tokens
}

func ConstructTree(token lexerV2Tokens.LexNodeToken) {
}

func consume(char rune) {
}

// Switch to a state, but when attempting to consume the next input character,
// provide the current input character instead.
func reconsume(char rune) {
}
