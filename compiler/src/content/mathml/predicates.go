package mathML

import "strings"

func IsMathMlFile(filename string) bool {
	return strings.HasSuffix(filename, ".mathml") || strings.HasSuffix(filename, ".mml")
}
