package markdown

import "strings"

func IsMarkdownFile(filePath string) bool {
	return strings.HasSuffix(filePath, ".md")
}
