package lexer

import "regexp"

type LexerToken *regexp.Regexp

func token(expression string) LexerToken {
	r, err := regexp.Compile(expression)
	if err != nil {
		panic(err)
	}

	return r
}
