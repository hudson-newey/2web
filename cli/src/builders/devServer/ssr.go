package devserver

import (
	"github.com/hudson-newey/2web-cli/src/shell"
)

func serveSsr() {
	shell.ExecuteCommand("./server/ssr.ts")
}
