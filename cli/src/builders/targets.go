package builders

import "os"

// Extracts the directory/file location that builders (vite, biome, etc...)
// should target.
// This is done by passing a third argument to the builder commands.
//
// 2web serve [path]
// E.g. "2web serve ." should serve the current directory.
func entryTarget(args []string) string {
	defaultEntryTarget := "./src/"

	if len(args) > 2 {
		return args[2]
	}

	return defaultEntryTarget
}

func hasSsrTarget() bool {
	_, err := os.Stat("./server/ssr.ts")
	return err == nil
}
