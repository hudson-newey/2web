package devserver

import "github.com/hudson-newey/2web/_shared/shell"

func serveSsr() {
	shell.ExecuteSync("./server/ssr.ts")
}
