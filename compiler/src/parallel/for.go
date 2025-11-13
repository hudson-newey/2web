package parallel

import (
	"hudson-newey/2web/src/cli"
	"sync"
)

// A for loop that runs each iteration in parallel.
func ForEach[T any](items []T, fn func(T)) {
	serial := *cli.GetArgs().Serial
	if serial {
		serialForEach(items, fn)
	} else {
		parallelForEach(items, fn)
	}
}

func serialForEach[T any](items []T, fn func(T)) {
	for _, value := range items {
		fn(value)
	}
}

func parallelForEach[T any](items []T, fn func(T)) {
	waitGroup := sync.WaitGroup{}

	// Because this blocks the main thread until all iterations are complete, we
	// can re-use the main thread for the last iteration to save on a goroutine.
	// This is useful because creating a goroutine has some overhead which we can
	// avoid for one iteration.
	synchronizedIndex := len(items) - 1

	for i, value := range items {
		waitGroup.Add(1)

		if i < synchronizedIndex {
			go func(v T) {
				defer waitGroup.Done()
				fn(v)
			}(value)
			continue
		}

		fn(value)
		waitGroup.Done()
	}

	waitGroup.Wait()
}
