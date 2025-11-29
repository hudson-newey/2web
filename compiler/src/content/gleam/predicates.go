package gleam

import "strings"

func IsGleamFile(path string) bool {
	return strings.HasSuffix(path, ".gleam")
}
