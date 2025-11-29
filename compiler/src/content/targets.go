package content

import (
	"strings"
)

// <name>.js.<language>
func IsBrowserTarget(path string) bool {
	return target(path) == "js" || target(path) == "browser"
}

// <name>.ssr.<language>
func IsServerTarget(path string) bool {
	return target(path) == "ssr" || target(path) == "server"
}

func IsWasmTarget(path string) bool {
	return target(path) == "wasm"
}

func target(path string) string {
	splits := strings.Split(path, ".")
	if len(splits) < 2 {
		return ""
	}

	target := splits[len(splits)-2]

	return target
}
