package devserver

import (
	"github.com/hudson-newey/2web-cli/src/ssr"
)

func ServeSolution(args []string) {
	if ssr.HasSsrTarget() {
		serveSsr()
	}

	serveSpa(args)
}
