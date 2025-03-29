package lexer

import "regexp"

func Token(expression string) *regexp.Regexp {
	r, err := regexp.Compile(expression)
	if err != nil {
		panic(err)
	}

	return r
}
