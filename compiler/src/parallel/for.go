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

	for _, value := range items {
		waitGroup.Add(1)

		go func(v T) {
			defer waitGroup.Done()
			fn(v)
		}(value)
	}

	waitGroup.Wait()
}
