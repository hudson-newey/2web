package document

import (
	"strings"
)

type InjectableLevel int

const (
	// Add to the start of content
	Leading InjectableLevel = iota

	// At the end of the HTML document
	Html

	// After the first <head> tag
	HeadTop
	// Before the last </head> tag
	Head

	// After the first <body> tag
	BodyTop
	// before the last </body> tag
	Body

	// Add to the end of content
	Trailing
)

// injects content into the bottom of the document while maintaining correct
// HTML structure
func InjectContent(content string, injectable string, level InjectableLevel) string {
	if level == Leading {
		return injectable + content
	} else if level == Trailing {
		return content + injectable
	}

	targetSelector := ""
	switch level {
	case Html:
		targetSelector = "</html>"

	case HeadTop:
		targetSelector = "<head>"
	case Head:
		targetSelector = "</head>"

	case BodyTop:
		targetSelector = "<body>"
	case Body:
		targetSelector = "</body>"
	}

	if level == HeadTop || level == BodyTop {
		return strings.Replace(content, targetSelector, targetSelector+injectable, 1)
	}

	return strings.Replace(content, targetSelector, injectable+targetSelector, 1)
}
