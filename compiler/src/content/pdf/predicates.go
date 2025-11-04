package pdf

import "strings"

func IsPdfFile(filePath string) bool {
	return strings.HasSuffix(filePath, ".pdf")
}
