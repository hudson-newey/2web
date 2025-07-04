package lexerV2

import (
	lexerV2State "hudson-newey/2web/src/compiler/v2/state"
	lexerV2Tokens "hudson-newey/2web/src/compiler/v2/tokens"
)

// https://html.spec.whatwg.org/multipage/parsing.html#overview-of-the-parsing-model
// https://html.spec.whatwg.org/multipage/parsing.html#tokenization
var scriptNestingLevel = 0
var parserPauseFlag = false
var currentState lexerV2State.State = lexerV2State.Data

func Tokenize(inputData string) []lexerV2Tokens.LexNodeToken {
	for _, char := range inputData {
	}
}

func ConstructTree(token lexerV2Tokens.LexNodeToken) {
}

func consume(char string) {
}

// Switch to a state, but when attempting to consume the next input character,
// provide the current input character instead.
func reconsume(char string) {
}
