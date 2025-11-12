package parallel

import (
	"hudson-newey/2web/src/cli"
	"sync"
)

// TODO: Replace this with sync.Go when I update to Go 1.25
// Run multiple functions in parallel and wait for all to complete.
func Run(fns ...func()) {
	serial := *cli.GetArgs().Serial
	if serial {
		serialRun(fns...)
	} else {
		parallelRun(fns...)
	}
}

func serialRun(fns ...func()) {
	for _, fn := range fns {
		fn()
	}
}

func parallelRun(fns ...func()) {
	waitGroup := sync.WaitGroup{}

	for _, fn := range fns {
		waitGroup.Add(1)

		go func(f func()) {
			defer waitGroup.Done()
			f()
		}(fn)
	}

	waitGroup.Wait()
}
