package lexer

import (
	"bufio"
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
	"io"
	"unicode"
)

type Lexer struct {
	pos    Position
	reader *bufio.Reader

	State LexFunc
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		pos:    Position{Row: 1, Col: 0},
		reader: bufio.NewReader(reader),
		State:  textLexer,
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

// This lexer is heavily inspired by arron raff's blog post
// https://www.aaronraff.dev/blog/how-to-write-a-lexer-in-go
func (model *Lexer) lex() V2LexNode {
	node, state := model.State(model)

	model.State = state

	return node
}

func textLexer(model *Lexer) (V2LexNode, LexFunc) {
	for {
		readerChar, _, err := model.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return V2LexNode{Pos: model.pos, Token: lexerTokens.EOF, Content: ""}, textLexer
			}

			panic(err)
		}

		switch readerChar {
		case '\n':
			model.lineFeed()
		case '<':
			return V2LexNode{Pos: model.pos, Token: lexerTokens.LessAngle, Content: "<"}, textLexer
		case '>':
			return V2LexNode{Pos: model.pos, Token: lexerTokens.GreaterAngle, Content: ">"}, textLexer
		case '/':
			return V2LexNode{Pos: model.pos, Token: lexerTokens.Slash, Content: "/"}, textLexer
		case '\'':
			return V2LexNode{Pos: model.pos, Token: lexerTokens.QuoteSingle, Content: "'"}, textLexer
		case '"':
			return V2LexNode{Pos: model.pos, Token: lexerTokens.QuoteDouble, Content: "\""}, textLexer
		case '@':
			return V2LexNode{Pos: model.pos, Token: lexerTokens.AtSymbol, Content: "@"}, textLexer
		case '*':
			return V2LexNode{Pos: model.pos, Token: lexerTokens.Star, Content: "*"}, textLexer
		case '#':
			return V2LexNode{Pos: model.pos, Token: lexerTokens.Hash, Content: "#"}, textLexer
		default:
			if unicode.IsSpace(readerChar) {
				continue
			} else if unicode.IsLetter(readerChar) {
				startPos := model.pos
				model.backup()
				text := model.lexText()
				return V2LexNode{Pos: startPos, Token: lexerTokens.Text, Content: string(text)}, textLexer
			}
		}
	}
}

// Returns the column position to the start and increments the line number
func (model *Lexer) lineFeed() {
	model.pos.Row++
	model.pos.Col = 0
}

func (model *Lexer) backup() {
	if err := model.reader.UnreadRune(); err != nil {
		panic(err)
	}

	model.pos.Col--
}

// lexIdent scans the input until the end of an identifier and then returns the
// literal.
func (l *Lexer) lexText() string {
	var lit string
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// at the end of the identifier
				return lit
			}
		}

		l.pos.Col++
		if unicode.IsLetter(r) {
			lit = lit + string(r)
		} else {
			// scanned something not in the identifier
			l.backup()
			return lit
		}
	}
}
