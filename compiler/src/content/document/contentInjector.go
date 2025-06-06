package document

import (
	"strings"
)

type InjectableLevel int

const (
	Html InjectableLevel = iota
	Head
	Body
	Trailing
)

// injects content into the bottom of the document while maintaining correct
// HTML structure
func InjectContent(content string, injectable string, level InjectableLevel) string {
	if level == Trailing {
		return content + injectable
	}

	targetSelector := ""
	switch level {
	case Html:
		targetSelector = "</html>"
	case Head:
		targetSelector = "</head>"
	case Body:
		targetSelector = "</body>"
	}

	return strings.Replace(content, targetSelector, injectable+targetSelector, 1)
}
