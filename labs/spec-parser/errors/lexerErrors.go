package lexerV2Errors

type ParseError = int

// TODO: Break down specific errors
// https://html.spec.whatwg.org/multipage/parsing.html#parse-errors
const (
	UnknownError ParseError = iota
)
