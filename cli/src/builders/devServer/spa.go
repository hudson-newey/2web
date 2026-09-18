package devserver

import (
	"slices"

	"github.com/hudson-newey/2web-cli/src/builders"
	"github.com/hudson-newey/2web-cli/src/server"
)

func serveSpa(args []string) {
	inPath := builders.EntryTargets(args)[0]
	outPath := builders.OutputTarget(args)

	server.Run(inPath, outPath, server.Options{
		// TODO: Refactor to use go flags package
		WatchFiles: !slices.Contains(args, "--no-watch"),
		AutoReload: !slices.Contains(args, "--no-auto-reload"),
	})
}
