package lexer

import (
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
	"io"
	"unicode"
)

// When inside the first starting angle bracket (<) and up until (and including)
// the closing angle bracket (>).
func elementLexer(model *Lexer) (V2LexNode, LexFunc) {
	for {
		readerChar, _, err := model.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return V2LexNode{Pos: model.pos, Token: lexerTokens.EOF, Content: ""}, textLexer
			}

			panic(err)
		}

		stringifiedChar := string(readerChar)
		switch readerChar {
		case '\n':
			model.lineFeed()
		case '/':
			return V2LexNode{Pos: model.pos, Token: lexerTokens.Slash, Content: stringifiedChar}, elementLexer
		case '\'':
			return V2LexNode{Pos: model.pos, Token: lexerTokens.QuoteSingle, Content: stringifiedChar}, elementLexer
		case '"':
			return V2LexNode{Pos: model.pos, Token: lexerTokens.QuoteDouble, Content: stringifiedChar}, elementLexer
		case '@':
			return V2LexNode{Pos: model.pos, Token: lexerTokens.AtSymbol, Content: stringifiedChar}, elementLexer
		case '*':
			return V2LexNode{Pos: model.pos, Token: lexerTokens.Star, Content: stringifiedChar}, elementLexer
		case '#':
			return V2LexNode{Pos: model.pos, Token: lexerTokens.Hash, Content: stringifiedChar}, elementLexer
		case '!':
			return V2LexNode{Pos: model.pos, Token: lexerTokens.Exclamation, Content: stringifiedChar}, elementLexer
		case '=':
			return V2LexNode{Pos: model.pos, Token: lexerTokens.Equals, Content: stringifiedChar}, elementLexer
		case '>':
			return V2LexNode{Pos: model.pos, Token: lexerTokens.GreaterAngle, Content: stringifiedChar}, textLexer
		default:
			if unicode.IsSpace(readerChar) {
				continue
			} else if unicode.IsLetter(readerChar) {
				startPos := model.pos
				model.backup()
				text := model.lexText()
				return V2LexNode{Pos: startPos, Token: lexerTokens.Text, Content: string(text)}, elementLexer
			}
		}
	}
}
