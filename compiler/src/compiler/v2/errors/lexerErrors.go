package lexerV2Errors

type parseError = int

// TODO: Break down specific errors
// https://html.spec.whatwg.org/multipage/parsing.html#parse-errors
const (
	UnknownError parseError = iota
)
