package builder

import (
	"log"
	"os"
	"sync"

	"hudson-newey/2web/src/builder/cache"
	"hudson-newey/2web/src/filesystem"
)

// pendingCacheEntry is an asset that has been compiled and written, but whose
// output has not yet been confirmed as flushed to disk.
type pendingCacheEntry struct {
	inputPath  string
	outputPath string
}

var pendingCacheMutex sync.Mutex
var pendingCacheEntries []pendingCacheEntry

// queueCacheAsset registers an asset to be added to the build cache once its
// output file has been confirmed as written to disk.
//
// Page compilation only enqueues file writes (see filesystem.WriteFile), so
// caching an asset before filesystem.WaitFileWriter has drained the write queue
// could record a cache entry for a file that was never written. A subsequent
// build would then trust the cache and never write the missing file.
func queueCacheAsset(inputPath string, outputPath string) {
	pendingCacheMutex.Lock()
	defer pendingCacheMutex.Unlock()

	pendingCacheEntries = append(pendingCacheEntries, pendingCacheEntry{
		inputPath,
		outputPath,
	})
}

// flushCacheEntries adds all queued assets to the build cache.
//
// It must only be called after filesystem.WaitFileWriter has confirmed that all
// queued file writes have been flushed to disk. Assets whose output write
// failed are skipped so that they are recompiled on the next build.
func flushCacheEntries() {
	pendingCacheMutex.Lock()
	entries := pendingCacheEntries
	pendingCacheEntries = nil
	pendingCacheMutex.Unlock()

	if len(entries) == 0 {
		return
	}

	failedWrites := make(map[string]bool)
	for _, path := range filesystem.FailedWrites() {
		failedWrites[path] = true
	}

	for _, entry := range entries {
		if failedWrites[entry.outputPath] {
			log.Printf("Skipping build cache for %s because its output could not be written", entry.outputPath)
			continue
		}

		if _, err := os.Stat(entry.outputPath); err != nil {
			log.Printf("Skipping build cache for %s because its output does not exist", entry.outputPath)
			continue
		}

		cache.CacheAsset(entry.inputPath, entry.outputPath)
	}
}
