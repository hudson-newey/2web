package lexer

import "regexp"

type LexToken *regexp.Regexp

func token(expression string) LexToken {
	r, err := regexp.Compile(expression)
	if err != nil {
		panic(err)
	}

	return r
}
