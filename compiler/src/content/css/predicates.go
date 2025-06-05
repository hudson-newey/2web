package css

import "strings"

func IsCssFile(filePath string) bool {
	return strings.HasSuffix(filePath, ".css")
}
