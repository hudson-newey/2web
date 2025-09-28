package html

import (
	"fmt"
	"strings"
	"unicode"
)

type headElement = string
type bodyElement = string

func ExpandPartial(content string) string {
	if isFullDocument(content) {
		return content
	}

	headContent, bodyContent := splitContent(content)

	// We automatically set the charset and viewport so that the developer doesn't
	// have to worry about this boilerplate.
	// TODO: This should be refactored when we introduce layouts.
	resultContent := fmt.Sprintf(
		`<!doctype html>
	<html>
		<head>
			<meta charset="UTF-8" />
			<meta name="viewport" content="width=device-width, initial-scale=1.0" />
			%s
		</head>

		<body>
			%s
		</body>
	</html>
`, headContent, bodyContent)

	return resultContent
}

func isFullDocument(content string) bool {
	return strings.HasPrefix(content, "<!")
}

func splitContent(content string) (headElement, bodyElement) {
	headElementTags := []string{"<meta", "<title", "<link", "<base"}

	headContent := ""
	bodyContent := ""

	lines := strings.SplitSeq(content, "\n")
	for line := range lines {
		// TODO: This doesn't cover the edge case where there are two elements
		// inside a single line.
		// E.g. <title>test</title><meta charset="UTF-8" />
		foundHeadElement := false
		for _, headSelector := range headElementTags {
			trimmedLine := strings.TrimLeftFunc(line, unicode.IsSpace)
			if strings.HasPrefix(trimmedLine, headSelector) {
				foundHeadElement = true
				break
			}
		}

		// Because we originally split the content by new lines, we need to re-add
		// the new lines to preserve the original document structure.
		if foundHeadElement {
			headContent += line + "\n"
		} else {
			bodyContent += line + "\n"
		}
	}

	return headContent, bodyContent
}
