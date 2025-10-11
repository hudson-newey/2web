package xslt

import "strings"

func IsXsltFile(filePath string) bool {
	return strings.HasSuffix(filePath, ".xslt") ||
		strings.HasSuffix(filePath, ".xsl")
}
