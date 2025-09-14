package lexer

import (
	"hudson-newey/2web/src/compiler/2-lexer/states"
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
	"io"
	"unicode"
)

type LexFunc func(*Lexer) (V2LexNode, LexFunc)

func lexerFactory(lexMap lexDefMap, state states.LexState) LexFunc {
	return func(lexerModel *Lexer) (V2LexNode, LexFunc) {
		for {
			matchingLexNode, nextState := lexMap.matching(lexerModel, state)
			if nextState != nil {
				return matchingLexNode, nextState
			}

			// Some systems use \r\n for new lines, but because we ignore \r, it will
			// work fine.
			readerChar, _, err := lexerModel.nextChar()
			if readerChar == '\n' {
				lexerModel.lineFeed()
			}

			if err != nil {
				if err == io.EOF {
					position := Position{
						Row: lexerModel.Pos.Row,
						Col: lexerModel.Pos.Col,
					}

					// We still return a textLexer state here so that if a file is partially
					// corrupted and contains an EOF before the real end, we can recover
					// gracefully.
					lexNode := V2LexNode{
						Pos:     position,
						Token:   lexerTokens.EOF,
						Content: "",
						State:   states.SourceText,
					}

					return lexNode, lexerModel.State
				}

				panic(err)
			}

			if unicode.IsSpace(readerChar) {
				continue
			} else if unicode.IsLetter(readerChar) || unicode.IsDigit(readerChar) {
				startPos := lexerModel.Pos
				lexerModel.backup(1)
				text := lexerModel.lexLiteral(lexMap)

				// There are different types of text depending on what context we are in
				// Sometimes it can be external source code.
				tokenMap := map[states.LexState]lexerTokens.LexToken{
					states.ScriptSource: lexerTokens.ScriptSource,
					states.StyleSource:  lexerTokens.ScriptSource,
					states.TextContent:  lexerTokens.TextContent,
				}

				token, exists := tokenMap[state]
				if !exists {
					token = lexerTokens.TextContent
				}

				position := Position{
					Row: startPos.Row,
					Col: startPos.Col,
				}

				lexNode := V2LexNode{
					Pos:     position,
					Token:   token,
					State:   state,
					Content: string(text),
				}

				return lexNode, lexerModel.State
			}
		}
	}
}
