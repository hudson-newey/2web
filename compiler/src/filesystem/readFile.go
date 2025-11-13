package filesystem

import (
	"os"
	"sync"
)

var (
	mu       sync.Mutex
	cache    = make(map[string][]byte)
	inFlight = make(map[string]*call)
)

type call struct {
	wg  sync.WaitGroup
	res []byte
	err error
}

// A thread-safe file read operation.
// If multiple reads are issued in parallel while one is already in progress,
// they will use the result of the first read instead of issuing multiple reads.
// Additionally, the files contents is cached for future reads after the first
// read.
func ReadFile(path string) ([]byte, error) {
	mu.Lock()
	// Return cached contents if present.
	if data, ok := cache[path]; ok {
		out := append([]byte(nil), data...)
		mu.Unlock()
		return out, nil
	}

	// If a read is already in-flight, wait for it.
	if c, ok := inFlight[path]; ok {
		mu.Unlock()
		c.wg.Wait()
		if c.err != nil {
			return nil, c.err
		}
		return append([]byte(nil), c.res...), nil
	}

	// Start a new in-flight call.
	c := &call{}
	c.wg.Add(1)
	inFlight[path] = c
	mu.Unlock()

	// Perform the actual read.
	data, err := os.ReadFile(path)

	// Cache on success.
	if err == nil {
		mu.Lock()
		cache[path] = append([]byte(nil), data...)
		mu.Unlock()
	}

	// Publish result to waiters and remove in-flight marker.
	c.res = append([]byte(nil), data...)
	c.err = err
	c.wg.Done()

	mu.Lock()
	delete(inFlight, path)
	mu.Unlock()

	if err != nil {
		return nil, err
	}
	return append([]byte(nil), data...), nil
}
