package txt

import "strings"

func IsTxtFile(filePath string) bool {
	return strings.HasSuffix(filePath, ".txt")
}
