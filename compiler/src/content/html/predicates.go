package html

import "strings"

func IsHtmlFile(filePath string) bool {
	return strings.HasSuffix(filePath, ".html")
}

func IsPage(filePath string) bool {
	return IsHtmlFile(filePath) && !IsComponent(filePath)
}

func IsComponent(filePath string) bool {
	return strings.HasSuffix(filePath, ".component.html")
}
