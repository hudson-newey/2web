package c

import "strings"

func IsCFile(filename string) bool {
	return strings.HasSuffix(filename, ".c")
}
