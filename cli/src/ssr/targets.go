package ssr

import "os"

const SsrTargetDir string = "./server/"

// The directory that the compiler emits compiled server routes into.
const ServerOutputDir string = "./dist-server/"

// The route manifest that the compiler writes into the server output
// directory. It is consumed by the server runtime to mount the compiled
// handlers.
const ServerRouteManifest string = ServerOutputDir + "routes.json"

func HasSsrTarget() bool {
	_, err := os.Stat(SsrTargetDir + "ssr.ts")
	return err == nil
}

// HasServerRoutes returns whether the solution has compiled server routes
// (i.e. the compiler has produced a server output directory with a route
// manifest).
func HasServerRoutes() bool {
	_, err := os.Stat(ServerRouteManifest)
	return err == nil
}
