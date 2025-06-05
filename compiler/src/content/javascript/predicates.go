package javascript

import "strings"

func IsJsFile(filePath string) bool {
	return (strings.HasSuffix(filePath, ".js") ||
		strings.HasSuffix(filePath, ".mjs") ||
		strings.HasSuffix(filePath, ".cjs"))
}
