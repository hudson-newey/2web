package lexer

import (
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
	"io"
	"unicode"
)

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
			return V2LexNode{Pos: model.pos, Token: lexerTokens.LessAngle, Content: "<"}, elementLexer
		case '>':
			return V2LexNode{Pos: model.pos, Token: lexerTokens.GreaterAngle, Content: ">"}, textLexer
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
