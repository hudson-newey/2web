package odt

import "strings"

func IsOdtFile(filePath string) bool {
	return strings.HasSuffix(filePath, ".odt")
}
