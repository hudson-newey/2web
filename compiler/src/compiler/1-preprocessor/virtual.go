package preprocessor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"hudson-newey/2web/src/content/document/documentErrors"
	"hudson-newey/2web/src/models"
)

// Compile time virtual functions are functions that the compiler evaluates
// while it compiles a page (they don't exist at runtime). They are imported
// in a <script compiled> block like a normal esm import:
//
//	import { $uid, $env, $readFile } from "@compiler/types/interop.types.ts";
//
// The import statement is removed (see expandComponents), and every call is
// replaced with its compile time value:
//
//   - $uid() expands to a unique (bare) incrementing id, e.g. "uid-3", which
//     is useful for aria attributes: 'aria-describedby="$uid()"'.
//   - $env("NAME") expands to a quoted string literal of the environment
//     variable's value. It is a compile time constant, not reactive.
//   - $readFile("path") expands to a quoted string literal of the file's
//     content (paths are resolved relative to the page). It is a compile
//     time constant, not reactive.
//
// The ids are deterministic per page (they increment in document order), so
// repeated builds compile identical output.

var (
	uidCallPattern      = regexp.MustCompile(`\$uid\(\)`)
	envCallPattern      = regexp.MustCompile(`\$env\(\s*"([^"]*)"\s*\)`)
	readFileCallPattern = regexp.MustCompile(`\$readFile\(\s*"([^"]*)"\s*\)`)
)

// expandVirtualFunctions replaces every compile time virtual function call in
// the page with the value the compiler evaluates for it.
func expandVirtualFunctions(filePath string, content string) string {
	uid := 0

	content = uidCallPattern.ReplaceAllStringFunc(content, func(string) string {
		uid++
		return "uid-" + strconv.Itoa(uid)
	})

	importDirectory := filepath.Dir(filePath)

	content = envCallPattern.ReplaceAllStringFunc(content, func(match string) string {
		name := envCallPattern.FindStringSubmatch(match)[1]

		return jsStringLiteral(os.Getenv(name))
	})

	content = readFileCallPattern.ReplaceAllStringFunc(content, func(match string) string {
		readPath := readFileCallPattern.FindStringSubmatch(match)[1]

		if !filepath.IsAbs(readPath) {
			readPath = filepath.Join(importDirectory, readPath)
		}

		content, err := os.ReadFile(readPath)
		if err != nil {
			readError := models.NewError(
				fmt.Sprintf("$readFile() could not read the file '%s': %v", readPath, err),
				filePath,
				models.Position{},
			)

			documentErrors.AddErrors(&readError)
			return `""`
		}

		return jsStringLiteral(string(content))
	})

	return content
}

// jsStringLiteral formats a string as a javascript string literal.
func jsStringLiteral(value string) string {
	encoded, err := json.Marshal(value)
	if err != nil {
		return `""`
	}

	return string(encoded)
}
