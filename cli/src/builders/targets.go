package builders

import "github.com/hudson-newey/2web-cli/src/ssr"

// Extracts the directory/file location that builders (vite, oxc, etc...)
// should target.
// This is done by passing a third argument to the builder commands.
//
// 2web serve [path]
// E.g. "2web serve ." should serve the current directory.
func EntryTargets(args []string) []string {
	if len(args) > 2 {
		return []string{args[2]}
	}

	defaultEntryTargets := []string{"./src/"}
	if ssr.HasSsrTarget() {
		defaultEntryTargets = append(defaultEntryTargets, ssr.SsrTargetDir)
	}

	return defaultEntryTargets
}

// TODO: Make this configurable via a cli flag
func OutputTarget(_args []string) string {
	return "./dist/"
}
