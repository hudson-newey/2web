package cache

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"hudson-newey/2web/src/constants"
	"hudson-newey/2web/src/filesystem"
)

// TestMain points the cache at a throwaway directory before anything in the
// package can resolve the environment (the resolved values are memoized), so
// that the database tests never touch the developer's real build cache.
func TestMain(m *testing.M) {
	cacheDir, err := os.MkdirTemp("", "2web-cache-test-")
	if err != nil {
		panic(err)
	}

	os.Setenv(constants.EnvCacheOverride, cacheDir+"/")

	code := m.Run()

	os.RemoveAll(cacheDir)
	os.Exit(code)
}

// resetCacheDB closes the shared connection and removes the database files so
// that each database test starts from an empty cache.
func resetCacheDB(t *testing.T) {
	t.Helper()

	CloseDBConnection()

	for _, suffix := range []string{"", "-wal", "-shm"} {
		os.Remove(dbLocation() + suffix)
	}
}

// insertCacheEntry writes a cache row directly, bypassing CacheAssets, so that
// tests can control the recorded last_touched time.
func insertCacheEntry(t *testing.T, key string, lastTouched int64) {
	t.Helper()

	conn := dbConnection()
	if conn == nil {
		t.Fatal("the cache database is unavailable")
	}

	if _, err := conn.Exec(
		`INSERT INTO `+buildCacheTableName+` (in_out_mod, last_touched) VALUES (?, ?);`,
		key, lastTouched,
	); err != nil {
		t.Fatalf("failed to insert cache entry: %v", err)
	}

	// The in memory snapshot is loaded when the connection is opened, so the
	// direct insert has to be added to it for lookups to see the entry.
	rememberCachedKey(key)
}

// entryCount returns the number of rows in the build cache.
func entryCount(t *testing.T) int64 {
	t.Helper()

	conn := dbConnection()
	if conn == nil {
		t.Fatal("the cache database is unavailable")
	}

	var count int64
	if err := conn.QueryRow(`SELECT COUNT(*) FROM ` + buildCacheTableName + `;`).Scan(&count); err != nil {
		t.Fatalf("failed to count cache entries: %v", err)
	}

	return count
}

// lastTouchedOf returns the recorded last_touched time of the given key.
func lastTouchedOf(t *testing.T, key string) int64 {
	t.Helper()

	conn := dbConnection()
	if conn == nil {
		t.Fatal("the cache database is unavailable")
	}

	var lastTouched int64
	if err := conn.QueryRow(
		`SELECT last_touched FROM `+buildCacheTableName+` WHERE in_out_mod = ?;`,
		key,
	).Scan(&lastTouched); err != nil {
		t.Fatalf("failed to read the last_touched time of %s: %v", key, err)
	}

	return lastTouched
}

func writeFile(t *testing.T, path string, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("failed to create directory for %s: %v", path, err)
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write %s: %v", path, err)
	}

	// File reads are cached for the lifetime of a compilation. Each key
	// computation in these tests simulates a fresh build, so the read cache
	// must be dropped to see the new file content.
	filesystem.ClearReadCache()
}

func keyFor(t *testing.T, inputPath string) string {
	t.Helper()

	key, err := CacheKey(inputPath, "/out/page.html")
	if err != nil {
		t.Fatalf("CacheKey returned an error: %v", err)
	}

	return key
}

// Regression tests for pages being served from cache after one of their
// compile time dependencies changed. The cache key is a fingerprint of the
// page's source and every file its compiled output depends on, so changing any
// of them must change the key.
func TestCacheKeyChangesWithPageContent(t *testing.T) {
	tempDir := t.TempDir()
	pagePath := filepath.Join(tempDir, "page.html")

	writeFile(t, pagePath, "<h1>v1</h1>")
	before := keyFor(t, pagePath)

	writeFile(t, pagePath, "<h1>v2</h1>")
	after := keyFor(t, pagePath)

	if before == after {
		t.Error("changing the page content did not change the cache key")
	}
}

func TestCacheKeyChangesWithLayoutChange(t *testing.T) {
	tempDir := t.TempDir()
	pagePath := filepath.Join(tempDir, "page.html")

	writeFile(t, pagePath, "<h1>Page</h1>")
	writeFile(t, filepath.Join(tempDir, "__layout.html"), "<html><slot></slot></html>")
	before := keyFor(t, pagePath)

	writeFile(t, filepath.Join(tempDir, "__layout.html"), "<html><body><slot></slot></body></html>")
	after := keyFor(t, pagePath)

	if before == after {
		t.Error("changing the shared __layout.html did not change the cache key")
	}
}

func TestCacheKeyChangesWithStyleSidecarChange(t *testing.T) {
	tempDir := t.TempDir()
	pagePath := filepath.Join(tempDir, "page.html")

	writeFile(t, pagePath, "<h1>Page</h1>")
	writeFile(t, filepath.Join(tempDir, "__style.css"), "h1 { color: red; }")
	before := keyFor(t, pagePath)

	writeFile(t, filepath.Join(tempDir, "__style.css"), "h1 { color: blue; }")
	after := keyFor(t, pagePath)

	if before == after {
		t.Error("changing the __style.css sidecar did not change the cache key")
	}
}

func TestCacheKeyChangesWithScriptSidecarChange(t *testing.T) {
	tempDir := t.TempDir()
	pagePath := filepath.Join(tempDir, "page.html")

	writeFile(t, pagePath, "<h1>Page</h1>")
	writeFile(t, filepath.Join(tempDir, "__script.ts"), "console.log('v1');")
	before := keyFor(t, pagePath)

	writeFile(t, filepath.Join(tempDir, "__script.ts"), "console.log('v2');")
	after := keyFor(t, pagePath)

	if before == after {
		t.Error("changing the __script.ts sidecar did not change the cache key")
	}
}

func TestCacheKeyChangesWithComponentImportChange(t *testing.T) {
	tempDir := t.TempDir()
	pagePath := filepath.Join(tempDir, "page.html")

	writeFile(t, pagePath, "<script compiled>\nimport Badge from \"components/badge.component.html\";\n</script>\n\n<Badge />")
	writeFile(t, filepath.Join(tempDir, "components", "badge.component.html"), "<span>v1</span>")
	before := keyFor(t, pagePath)

	writeFile(t, filepath.Join(tempDir, "components", "badge.component.html"), "<span>v2</span>")
	after := keyFor(t, pagePath)

	if before == after {
		t.Error("changing an imported component did not change the cache key")
	}
}

// Components can import other components, so a change to a transitive
// dependency must invalidate the page too.
func TestCacheKeyChangesWithTransitiveComponentChange(t *testing.T) {
	tempDir := t.TempDir()
	pagePath := filepath.Join(tempDir, "page.html")

	writeFile(t, pagePath, "<script compiled>\nimport Header from \"components/header.component.html\";\n</script>\n\n<Header />")
	writeFile(t, filepath.Join(tempDir, "components", "header.component.html"), "<script compiled>\nimport Title from \"./title.component.html\";\n</script>\n<Title />")
	writeFile(t, filepath.Join(tempDir, "components", "title.component.html"), "<title>v1</title>")
	before := keyFor(t, pagePath)

	writeFile(t, filepath.Join(tempDir, "components", "title.component.html"), "<title>v2</title>")
	after := keyFor(t, pagePath)

	if before == after {
		t.Error("changing a transitively imported component did not change the cache key")
	}
}

// ESM imports inside plain <script> tags are bundled by esbuild at compile
// time, so they are compile time dependencies as well.
func TestCacheKeyChangesWithEsmImportChange(t *testing.T) {
	tempDir := t.TempDir()
	pagePath := filepath.Join(tempDir, "page.html")

	writeFile(t, pagePath, "<script>\nimport { hello } from \"./scripts/esm.ts\";\n</script>")
	writeFile(t, filepath.Join(tempDir, "scripts", "esm.ts"), "export const hello = 'v1';")
	before := keyFor(t, pagePath)

	writeFile(t, filepath.Join(tempDir, "scripts", "esm.ts"), "export const hello = 'v2';")
	after := keyFor(t, pagePath)

	if before == after {
		t.Error("changing an ESM import did not change the cache key")
	}
}

// Dependencies that the compiler doesn't read must not invalidate the page.
func TestCacheKeyIgnoresUnrelatedFiles(t *testing.T) {
	tempDir := t.TempDir()
	pagePath := filepath.Join(tempDir, "page.html")

	writeFile(t, pagePath, "<h1>Page</h1>")
	writeFile(t, filepath.Join(tempDir, "__layout.html"), "<slot></slot>")
	before := keyFor(t, pagePath)

	writeFile(t, filepath.Join(tempDir, "unrelated.html"), "<p>not a dependency</p>")
	after := keyFor(t, pagePath)

	if before != after {
		t.Error("changing an unrelated file changed the cache key")
	}
}

func TestCacheKeyIsDeterministic(t *testing.T) {
	tempDir := t.TempDir()
	pagePath := filepath.Join(tempDir, "page.html")

	writeFile(t, pagePath, "<h1>Page</h1>")
	writeFile(t, filepath.Join(tempDir, "__layout.html"), "<slot></slot>")
	writeFile(t, filepath.Join(tempDir, "components", "badge.component.html"), "<span>v1</span>")

	first := keyFor(t, pagePath)
	second := keyFor(t, pagePath)

	if first != second {
		t.Error("CacheKey is not deterministic for identical inputs")
	}
}

func TestCacheKeyIncludesOutputPath(t *testing.T) {
	tempDir := t.TempDir()
	pagePath := filepath.Join(tempDir, "page.html")

	writeFile(t, pagePath, "<h1>Page</h1>")

	first, err := CacheKey(pagePath, "/out/one.html")
	if err != nil {
		t.Fatalf("CacheKey returned an error: %v", err)
	}

	second, err := CacheKey(pagePath, "/out/two.html")
	if err != nil {
		t.Fatalf("CacheKey returned an error: %v", err)
	}

	if first == second {
		t.Error("cache keys for different output paths should differ")
	}

	if !strings.HasPrefix(first, pagePath+"->"+"\""+"/out/one.html") && !strings.Contains(first, "/out/one.html") {
		t.Errorf("cache key %q does not reference its output path", first)
	}
}

// A page that can't be read can't be fingerprinted, so it must be treated as
// uncached (and uncacheable) instead of panicking or serving stale output.
func TestCacheKeyFailsForMissingPage(t *testing.T) {
	missingPage := filepath.Join(t.TempDir(), "does-not-exist.html")

	if _, err := CacheKey(missingPage, "/out/page.html"); err == nil {
		t.Error("CacheKey should return an error for a missing page file")
	}
}

// CacheAssets records entries with the build's start time, and every entry
// written during one build shares that exact timestamp (the clock is only
// queried once per build, see BeginBuild).
func TestCacheAssetsRecordsBuildTimestamp(t *testing.T) {
	resetCacheDB(t)

	buildStart := time.Now()
	BeginBuild(buildStart)

	pagePath := filepath.Join(t.TempDir(), "page.html")
	writeFile(t, pagePath, "<h1>Page</h1>")

	first := keyFor(t, pagePath)
	CacheAssets([]CacheRecord{{InputPath: pagePath, OutputPath: "/out/page.html", Key: first}})

	// Simulate wall clock time passing within the same build.
	time.Sleep(1100 * time.Millisecond)

	second := keyFor(t, pagePath)
	CacheAssets([]CacheRecord{{InputPath: pagePath, OutputPath: "/out/other.html", Key: second}})

	if touched := lastTouchedOf(t, first); touched != buildStart.Unix() {
		t.Errorf("expected the first entry to be touched at the build start, got %d", touched)
	}

	if touched := lastTouchedOf(t, second); touched != buildStart.Unix() {
		t.Errorf("expected the second entry to share the build's timestamp, got %d", touched)
	}
}

// A cache hit is a use of an entry: flushing the touches of a build updates
// the entry's last_touched time to the build's timestamp.
func TestFlushTouchesUpdatesHitEntries(t *testing.T) {
	resetCacheDB(t)

	buildStart := time.Now()
	BeginBuild(buildStart)

	hitPath := filepath.Join(t.TempDir(), "hit.html")
	untouchedPath := filepath.Join(t.TempDir(), "untouched.html")
	writeFile(t, hitPath, "<h1>hit</h1>")
	writeFile(t, untouchedPath, "<h1>untouched</h1>")

	hitKey := "stale->" + hitPath + "@key"
	untouchedKey := "stale->" + untouchedPath + "@key"

	insertCacheEntry(t, hitKey, buildStart.Unix()-3600)
	insertCacheEntry(t, untouchedKey, buildStart.Unix()-3600)

	if !IsCached("stale", hitPath, hitKey) {
		t.Fatal("expected the hit key to be cached")
	}

	FlushTouches()

	if touched := lastTouchedOf(t, hitKey); touched != buildStart.Unix() {
		t.Errorf("expected the hit entry to be touched at the build start, got %d", touched)
	}

	if touched := lastTouchedOf(t, untouchedKey); touched == buildStart.Unix() {
		t.Error("the untouched entry should keep its old last_touched time")
	}
}

// A cache under the size limit is never vacuumed, whatever the ages of its
// entries.
func TestVacuumKeepsCacheUnderLimit(t *testing.T) {
	resetCacheDB(t)

	BeginBuild(time.Now())
	previousMax := maxCacheSize
	maxCacheSize = 1 << 30 // 1 GiB: everything fits
	defer func() { maxCacheSize = previousMax }()

	old := buildTime() - 400*24*3600
	insertCacheEntry(t, "old->/out/a.html@k", old)
	insertCacheEntry(t, "old->/out/b.html@k", old)

	Vacuum()

	if entryCount(t) != 2 {
		t.Errorf("expected the vacuum to be a no-op for a cache under the size limit, got %d entries", entryCount(t))
	}
}

// Once the cache is over the size limit, entries that builds haven't touched
// within the retention window are removed first, and recently touched entries
// survive (the cache is cut back down around them).
func TestVacuumRemovesStaleEntriesFirst(t *testing.T) {
	resetCacheDB(t)

	BeginBuild(time.Now())
	previousMax := maxCacheSize
	maxCacheSize = 64 << 10
	defer func() { maxCacheSize = previousMax }()

	conn := dbConnection()
	if conn == nil {
		t.Fatal("the cache database is unavailable")
	}

	now := buildTime()

	// Two entries that no build has touched for months, and enough recently
	// touched entries to keep the cache over the size limit without them.
	if _, err := conn.Exec(
		`INSERT INTO `+buildCacheTableName+` (in_out_mod, last_touched) VALUES (?, ?), (?, ?);`,
		"stale->/out/a.html@k", now-90*24*3600,
		"stale->/out/b.html@k", now-60*24*3600,
	); err != nil {
		t.Fatal(err)
	}

	tx, err := conn.Begin()
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 300; i++ {
		key := "fresh->/out/page-" + strconv.Itoa(i) + ".html@" + strings.Repeat("f", 100)
		if _, err := tx.Exec(
			`INSERT INTO `+buildCacheTableName+` (in_out_mod, last_touched) VALUES (?, ?);`,
			key, now,
		); err != nil {
			t.Fatal(err)
		}
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	Vacuum()

	if hasCachedKey("stale->/out/a.html@k") || hasCachedKey("stale->/out/b.html@k") {
		t.Error("expected the stale entries to be removed")
	}

	if entryCount(t) < 20 {
		t.Errorf("expected most recently touched entries to survive the vacuum, got %d", entryCount(t))
	}

	if after := cacheSize(); after > maxCacheSize {
		t.Errorf("expected the vacuum to cut the cache below %d bytes, got %d bytes", maxCacheSize, after)
	}
}

// A cache that is still over the limit after the stale entries are removed is
// cut back to the low water mark by evicting the oldest entries first.
func TestVacuumEvictsOldestWhenStillOverLimit(t *testing.T) {
	resetCacheDB(t)

	BeginBuild(time.Now())

	previousMax := maxCacheSize
	previousRetention := retentionSeconds
	maxCacheSize = 128 << 10
	// A retention window that no entry can fall inside, so that the age stage
	// never fires and the eviction stage does all of the work.
	retentionSeconds = 100 * 365 * 24 * 3600
	defer func() {
		maxCacheSize = previousMax
		retentionSeconds = previousRetention
	}()

	// 512 entries of ~1 KiB of key text each: over the 128 KiB limit, with
	// evenly spaced last_touched times.
	conn := dbConnection()
	if conn == nil {
		t.Fatal("the cache database is unavailable")
	}

	tx, err := conn.Begin()
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 512; i++ {
		key := strings.Repeat("k", 1000) + "->/out/" + string(rune('a'+i%26)) + "@" + time.Unix(int64(i)*60, 0).Format("150405")
		if _, err := tx.Exec(
			`INSERT INTO `+buildCacheTableName+` (in_out_mod, last_touched) VALUES (?, ?);`,
			key, int64(i)*60,
		); err != nil {
			t.Fatal(err)
		}
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	Vacuum()

	after := cacheSize()
	if after > maxCacheSize {
		t.Errorf("expected the vacuum to cut the cache below %d bytes, got %d bytes", maxCacheSize, after)
	}

	if entryCount(t) >= 512 {
		t.Errorf("expected old entries to be evicted, still have %d entries", entryCount(t))
	}
}

// Entries with no last_touched time (e.g. rows migrated from a cache created
// before the column existed) count as ancient and are vacuumed first.
func TestVacuumRemovesUntouchedLegacyEntries(t *testing.T) {
	resetCacheDB(t)

	BeginBuild(time.Now())
	previousMax := maxCacheSize
	maxCacheSize = 1 << 20
	defer func() { maxCacheSize = previousMax }()

	insertCacheEntry(t, "legacy->/out/a.html@k", 0)
	insertCacheEntry(t, "fresh->/out/b.html@k", buildTime())

	Vacuum()

	if !hasCachedKey("legacy->/out/a.html@k") {
		t.Error("expected the legacy entry to be removed")
	}

	if !hasCachedKey("fresh->/out/b.html@k") {
		t.Error("expected the touched entry to survive")
	}
}
