package lexeme

const (
	// A special lexer token that can be used in grammars to match against any
	// lexer token.
	// If you use the ANYTHING token, make sure that you have a really good
	// exit condition after the ANYTHING to prevent capturing everything past
	// this token.
	//
	// This token is purposely not exported so that people can use the more
	// stable lexerTokens.CaptureUntil() function.
	CaptureUntil Lexeme = "CAPTURE_UNTIL"
)

// By using a function here, we provide a more stable interface to the anything
// token.
// It also makes the Anything() call stand out more in syntax highlighting.
func NewCaptureUntil() Lexeme {
	return newSpecialToken(CaptureUntil)
}

func IsSpecialToken(token Lexeme, match Lexeme) bool {
	return token == newSpecialToken(match)
}

const specialPrefix string = "__special__"

func newSpecialToken(readableStr Lexeme) Lexeme {
	return Lexeme(specialPrefix) + readableStr
}
