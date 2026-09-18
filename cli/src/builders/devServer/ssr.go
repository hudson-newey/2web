package devserver

import "github.com/hudson-newey/2web/_shared/shell"

func serveSsr() {
	shell.ExecuteSubProcess("./server/ssr.ts")
}
