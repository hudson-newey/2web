package lexerV2

// A buffer that can be used to store temporary state for lexing
// E.g. to store character reference state
// https://html.spec.whatwg.org/multipage/parsing.html#character-reference-state
var tempBuffer = ""

func pushToBuffer(char string) {
	tempBuffer = tempBuffer + char
}

func flushBuffer() string {
	bufferValue := tempBuffer
	tempBuffer = ""

	return bufferValue
}
