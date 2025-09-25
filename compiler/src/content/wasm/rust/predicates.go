package rust

import "strings"

func IsRustFile(filename string) bool {
	return strings.HasSuffix(filename, ".rs")
}
