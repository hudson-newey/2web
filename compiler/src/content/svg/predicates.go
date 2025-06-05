package svg

import "strings"

func IsSvgFile(filePath string) bool {
	return strings.HasSuffix(filePath, ".svg")
}
