package lexer

import (
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
	"hudson-newey/2web/src/compiler/io/reader"
	"hudson-newey/2web/src/models"
	"io"
	"strings"
)

type Lexer struct {
	Pos   *Position
	Input *reader.Reader
	State LexFunc

	// Errors records the errors produced while lexing (e.g. source read
	// failures). They are attributed to the file being lexed and are rendered
	// into the page's error overlay by the builder.
	Errors []*models.Error
}

// RecordError records an error produced while lexing.
func (model *Lexer) RecordError(message string, position Position) {
	errorModel := models.NewError(message, model.Input.FilePath, position)
	model.Errors = append(model.Errors, &errorModel)
}

func NewLexer(reader *reader.Reader) *Lexer {
	return &Lexer{
		Pos:   &Position{Row: 1, Col: 1},
		Input: reader,
		State: textLexer,
	}
}

func (model *Lexer) Execute() []*V2LexNode {
	representation := []*V2LexNode{}

	for {
		lexNode := model.lex()
		representation = append(representation, &lexNode)

		if lexNode.Token == lexeme.EOF {
			return representation
		}
	}
}

func (model *Lexer) nextChar() (char rune, size int, err error) {
	return model.Input.Reader.ReadRune()
}

func (model *Lexer) peek(length int) string {
	return string(model.peekBytes(length))
}

// peekBytes peeks the next `length` bytes without copying them.
//
// The returned slice aliases the read buffer and is only valid until the next
// read operation, which is enough for the matcher comparisons that look but
// don't consume.
func (model *Lexer) peekBytes(length int) []byte {
	bytes, err := model.Input.Reader.Peek(length)
	if err != nil && err != io.EOF {
		// A read failure (other than end of file) must not kill the build.
		model.RecordError("failed to read source: "+err.Error(), *model.Pos)
		return nil
	}

	return bytes
}

func (model *Lexer) skip(length int) {
	for range length {
		char, _, err := model.Input.Reader.ReadRune()
		if err != nil && err != io.EOF {
			model.RecordError("failed to read source: "+err.Error(), *model.Pos)
			return
		}

		if char == '\n' {
			model.lineFeed()
			continue
		}

		model.Pos.Col++
	}
}

func (model *Lexer) backup(length int) {
	for range length {
		if err := model.Input.Reader.UnreadRune(); err != nil {
			model.RecordError("failed to unread source: "+err.Error(), *model.Pos)
			return
		}

		model.Pos.Col--
	}
}

// This lexer is heavily inspired by arron raff's blog post
// https://www.aaronraff.dev/blog/how-to-write-a-lexer-in-go
func (model *Lexer) lex() V2LexNode {
	node, state := model.State(model)

	model.State = state

	return node
}

// Returns the column position to the start and increments the line number
func (model *Lexer) lineFeed() {
	model.Pos.Row++
	model.Pos.Col = 1
}

// lexIdent scans the input until the end of an identifier and then returns the
// literal source that was scanned up until the first lexer exit condition.
// lexLiteral captures text until one of the exit conditions matches at the
// current position.
//
// The literal is accumulated in a strings.Builder. Concatenating character by
// character with `literal += string(char)` would be quadratic in the length of
// the literal because every concatenation copies the whole string so far.
func (model *Lexer) lexLiteral(exitConditions *compiledMatchers) string {
	var literal strings.Builder

	for {
		if exitConditions.wouldMatchAt(model) {
			return literal.String()
		}

		nextChar, _, err := model.Input.Reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return literal.String()
			}
		}

		if nextChar == '\n' {
			// Track line feeds so that lexer positions stay accurate for
			// compiler errors. Without this, every position after the first
			// multi line text node would point at the wrong row.
			model.lineFeed()
		} else {
			model.Pos.Col++
		}

		literal.WriteRune(nextChar)
	}
}
