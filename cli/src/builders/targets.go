package builders

// Extracts the directory/file location that builders (vite, biome, etc...)
// should target.
// This is done by passing a third argument to the builder commands.
//
// 2web serve [path]
// E.g. "2web serve ." should serve the current directory.
func EntryTarget(args []string) string {
	defaultEntryTarget := "./src/"

	if len(args) > 2 {
		return args[2]
	}

	return defaultEntryTarget
}

// TODO: Make this configurable via a cli flag
func OutputTarget(_args []string) string {
	return "./dist/"
}
