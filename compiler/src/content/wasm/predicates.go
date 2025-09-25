package wasm

import "strings"

func IsWasmFile(filename string) bool {
	return strings.HasSuffix(filename, ".wasm")
}
