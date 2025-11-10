package parallel

import "sync"

// Run multiple functions in parallel and wait for all to complete.
func Run(fns ...func()) {
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
