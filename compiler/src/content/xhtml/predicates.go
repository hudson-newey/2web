package xhtml

import "strings"

func IsXhtmlFile(filePath string) bool {
	return strings.HasSuffix(filePath, ".xhtml")
}
