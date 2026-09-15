package models

// Position is a location in a source file.
//
// It lives in the models package (rather than the lexer) so that both the
// lexer and the error models can reference it without an import cycle.
type Position struct {
	Row int
	Col int
}

// The position of the very start of a file.
var StartingPosition = Position{Row: 0, Col: 0}
