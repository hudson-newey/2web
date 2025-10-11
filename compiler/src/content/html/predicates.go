package html

import "strings"

func IsHtmlFile(filePath string) bool {
	return strings.HasSuffix(filePath, ".html") ||
		strings.HasSuffix(filePath, ".htm")
}
