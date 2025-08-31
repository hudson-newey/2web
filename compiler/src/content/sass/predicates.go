package scss

import "strings"

// Note that we treat sass and scss files as both "sass" files because the
// godartsass transpiler is transparent about sass/scss distinctions.
func IsSassFile(filename string) bool {
	return strings.HasSuffix(filename, ".sass") || strings.HasSuffix(filename, ".scss")
}
