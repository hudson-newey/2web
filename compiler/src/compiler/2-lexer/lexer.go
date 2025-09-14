package lexer

import (
	"hudson-newey/2web/src/compiler/2-lexer/states"
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
	"hudson-newey/2web/src/compiler/io/reader"
	"io"
)

type Lexer struct {
	Pos   *Position
	Input *reader.Reader
	State LexFunc
}

func NewLexer(reader *reader.Reader) *Lexer {
	return &Lexer{
		Pos:   &Position{Row: 1, Col: 1},
		Input: reader,
		State: textLexer,
	}
}

func (model *Lexer) Execute() []V2LexNode {
	representation := []V2LexNode{}

	for {
		lexNode := model.lex()
		representation = append(representation, lexNode)

		if lexNode.Token == lexerTokens.EOF {
			return representation
		}
	}
}

func (model *Lexer) nextChar() (char rune, size int, err error) {
	return model.Input.Reader.ReadRune()
}

func (model *Lexer) peek(length int) string {
	bytes, err := model.Input.Reader.Peek(length)
	if err != nil && err != io.EOF {
		panic(err)
	}

	return string(bytes)
}

func (model *Lexer) skip(length int) {
	for range length {
		char, _, err := model.Input.Reader.ReadRune()
		if err != nil && err != io.EOF {
			panic(err)
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
			panic(err)
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
func (model *Lexer) lexLiteral(exitConditions lexDefMap) string {
	var literal string
	for {
		r, _, err := model.Input.Reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// at the end of the identifier
				return literal
			}
		}

		model.Pos.Col++

		literal = literal + string(r)

		_, nextState := exitConditions.matching(model, states.SourceText)
		if nextState != nil {
			// We've reached an exit condition
			// Back up one character and return the literal.
			// When the original lexer resumes, it will re-consume the character that
			// caused the exit condition.
			model.backup(1)
			return literal
		}
	}
}
