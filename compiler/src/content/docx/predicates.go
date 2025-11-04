package docx

import "strings"

func IsDocxFile(filePath string) bool {
	// TODO: We should add support for .doc files too
	return strings.HasSuffix(filePath, ".docx")
}
