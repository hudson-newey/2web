package latex

import "strings"

func IsLatexFile(filename string) bool {
	return strings.HasSuffix(filename, ".tex")
}
