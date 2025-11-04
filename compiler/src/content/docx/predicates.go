package docx

import "strings"

func IsDocxFile(filePath string) bool {
	return strings.HasSuffix(filePath, ".docx") ||
		strings.HasSuffix(filePath, ".doc")
}
