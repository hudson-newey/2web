package twoscript

import "strings"

func IsTwoScriptFile(filePath string) bool {
	return strings.HasSuffix(filePath, ".2s")
}
