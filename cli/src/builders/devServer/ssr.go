package devserver

import "github.com/hudson-newey/2web-cli/src/runner"

func serveSsr() {
	runner.ExecuteScript("./server/ssr.ts")
}
