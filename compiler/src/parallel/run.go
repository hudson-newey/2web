package parallel

import (
	"hudson-newey/2web/src/cli"
	"sync"
)

// TODO: Replace this with sync.Go when I update to Go 1.25
// Run multiple functions in parallel and wait for all to complete.
func Run(fns ...func()) {
	serial := cli.GetArgs().Serial
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

	// We re-use the main thread for the last function to save on a goroutine.
	// This is useful because creating a goroutine has some overhead which we can
	// avoid for one callback.
	synchronizedIndex := len(fns) - 1

	for i, fn := range fns {
		waitGroup.Add(1)

		if i < synchronizedIndex {
			go func(f func()) {
				defer waitGroup.Done()
				f()
			}(fn)
		} else {
			fn()
			waitGroup.Done()
		}
	}

	waitGroup.Wait()
}
