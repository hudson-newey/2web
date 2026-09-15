package lexeme

import "strings"

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

// The optional special token prefix is followed by the lexeme that it
// optionally matches. e.g. "__special__optional_TextContent".
const optionalPrefix string = specialPrefix + "optional_"

// The alternation special token prefix is followed by the lexemes that it can
// match, joined by the alternation separator.
// e.g. "__special__or_QuoteDouble|QuoteSingle".
const orPrefix string = specialPrefix + "or_"

// The separator used to join alternatives inside of an alternation special
// token.
const orSeparator string = "|"

// An optional token in a grammar definition matches zero or one occurrence of
// the wrapped lexeme.
//
// This is not a lexeme that the lexer can emit, it is a grammar matcher token
// that is created with NewOptional() and expanded by the grammar's Match
// method.
const Optional Lexeme = "OPTIONAL"

// An alternation token in a grammar definition matches exactly one occurrence
// of any of the wrapped lexemes.
//
// Like Optional, this is not a lexeme that the lexer can emit, it is a grammar
// matcher token that is created with NewOr() and expanded by the grammar's
// Match method.
const Or Lexeme = "OR"

// NewOptional wraps a lexeme so that a grammar definition matches zero or one
// occurrence of it.
//
// e.g. an optional TextContent token allows a grammar to tolerate whitespace
// between two tokens:
//
//	newDefinition(
//		lexeme.ControlFlowIfKeyword,
//		lexeme.NewOptional(lexeme.TextContent), // whitespace
//		lexeme.BracketOpen,
//	)
func NewOptional(def Lexeme) Lexeme {
	if IsOptionalToken(def) {
		// Optionally matching an optional token is still just optional.
		return def
	}

	return Lexeme(optionalPrefix) + def
}

func IsOptionalToken(token Lexeme) bool {
	return strings.HasPrefix(string(token), optionalPrefix)
}

// UnwrapOptional returns the lexeme that an optional token wraps.
func UnwrapOptional(token Lexeme) Lexeme {
	return Lexeme(strings.TrimPrefix(string(token), optionalPrefix))
}

// NewOr creates a grammar matcher token that matches exactly one occurrence of
// any one of the given lexemes.
//
// e.g. a grammar that accepts either quote style:
//
//	newDefinition(
//		lexeme.NewOr(lexeme.QuoteDouble, lexeme.QuoteSingle),
//		lexeme.TextContent,
//		lexeme.NewOr(lexeme.QuoteDouble, lexeme.QuoteSingle),
//	)
func NewOr(alternatives ...Lexeme) Lexeme {
	joined := strings.Join(alternativeStrings(alternatives), orSeparator)
	return Lexeme(orPrefix) + Lexeme(joined)
}

func IsOrToken(token Lexeme) bool {
	return strings.HasPrefix(string(token), orPrefix)
}

// UnwrapOr returns all lexemes that an alternation token can match.
func UnwrapOr(token Lexeme) []Lexeme {
	alternatives := strings.Split(strings.TrimPrefix(string(token), orPrefix), orSeparator)

	tokens := make([]Lexeme, len(alternatives))
	for i, alternative := range alternatives {
		tokens[i] = Lexeme(alternative)
	}

	return tokens
}

func alternativeStrings(alternatives []Lexeme) []string {
	alternativesAsStrings := make([]string, len(alternatives))
	for i, alternative := range alternatives {
		// Nested special tokens are not supported.
		alternativesAsStrings[i] = string(alternative)
	}

	return alternativesAsStrings
}

func newSpecialToken(readableStr Lexeme) Lexeme {
	return Lexeme(specialPrefix) + readableStr
}
