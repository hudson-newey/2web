package lexer

import (
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
	"io"
	"unicode"
)

func elementLexer(model *Lexer) (V2LexNode, LexFunc) {
	for {
		readerChar, _, err := model.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return V2LexNode{Pos: model.pos, Token: lexerTokens.EOF, Content: ""}, textLexer
			}

			panic(err)
		}

		switch readerChar {
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
