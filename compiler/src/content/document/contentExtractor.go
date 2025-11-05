package document

import "strings"

func ExtractTagContent(htmlContent string, tagName string) string {
	openTag := "<" + tagName + ">"
	closeTag := "</" + tagName + ">"

	startIndex := -1
	endIndex := -1

	startIndex = strings.Index(htmlContent, openTag)
	if startIndex == -1 {
		return ""
	}
	startIndex += len(openTag)
	endIndex = strings.Index(htmlContent, closeTag)
	if endIndex == -1 {
		return ""
	}

	return htmlContent[startIndex:endIndex]
}
