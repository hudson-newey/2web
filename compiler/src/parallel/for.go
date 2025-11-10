package parallel

import "sync"

// A for loop that runs each iteration in parallel.
func ForEach[T any](items []T, fn func(T)) {
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
