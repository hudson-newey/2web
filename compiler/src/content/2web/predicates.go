package twoWeb

import "strings"

func IsTwoWebFile(filePath string) bool {
	return strings.HasSuffix(filePath, ".2web")
}
