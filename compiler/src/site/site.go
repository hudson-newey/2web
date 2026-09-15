package site

import "sync"

// paths is mutated from the parallel build pool (BeforeEach is called by every
// page compilation worker), so all access must be synchronized.
var pathsMutex sync.Mutex
var paths []string

func GetSitePaths() []string {
	pathsMutex.Lock()
	defer pathsMutex.Unlock()

	// Return a copy so that callers can't race on the underlying slice while
	// page compilation workers are still appending to it.
	pathsCopy := make([]string, len(paths))
	copy(pathsCopy, paths)
	return pathsCopy
}

func registerSitePage(path string) {
	pathsMutex.Lock()
	defer pathsMutex.Unlock()

	paths = append(paths, path)
}
