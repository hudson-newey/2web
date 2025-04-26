package reactiveCompiler

import "strings"

func useDoubleQuotes(content string) bool {
	// we want to find out if we have to use single or double quotes
	// depending on what quote type the reducer uses
	firstDoubleQuoteLocation := strings.Index(content, "\"")
	firstSingleQuoteLocation := strings.Index(content, "'")

	if firstDoubleQuoteLocation == -1 {
		return true
	}
	if firstSingleQuoteLocation == -1 {
		return false
	}

	// We use greater than or equals to here so that the quotes default to
	// double. Because if the location of the first single and first double are
	// the same, then the reducer doesn't contain any quotes, so we can default
	// to double quotes
	return firstSingleQuoteLocation <= firstDoubleQuoteLocation
}
