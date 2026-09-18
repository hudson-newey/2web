package devserver

import "github.com/hudson-newey/2web/_shared/shell"

func serveSsr() {
	shell.ExecuteCommand("./server/ssr.ts")
}
