package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/hudson-newey/2web/_shared/logger"
	"hudson-newey/2web/src/cli"
	"hudson-newey/2web/src/filesystem"
)

// buildTimestamp is the time (in unix seconds) that every cache touch of the
// current build is recorded with.
//
// It is captured once at the start of compilation (see BeginBuild) so that the
// cache writes never need to query the clock themselves: every entry written
// or read during a build shares one timestamp, which is all the "last touched"
// resolution the vacuum needs (a build either touched an entry or it didn't).
var buildTimestamp int64 = 0

// BeginBuild records the build's timestamp for every cache touch made during
// this build. It must be called at the start of every compilation.
func BeginBuild(timestamp time.Time) {
	touchedKeysMutex.Lock()
	touchedKeys = nil
	touchedKeysMutex.Unlock()

	buildTimestamp = timestamp.Unix()
}

// buildTime returns the timestamp that cache touches of the current build are
// recorded with.
//
// Callers that use the cache outside of a build (no BeginBuild call) fall back
// to the current time.
func buildTime() int64 {
	if buildTimestamp == 0 {
		return time.Now().Unix()
	}

	return buildTimestamp
}

// touchedKeys collects the cache keys that were read (hit) during the current
// build. They are flushed into one batched last_touched update at the end of
// the build, so that a cache hit doesn't cost a write per page.
var touchedKeysMutex sync.Mutex
var touchedKeys []string

// touchKey records that the current build used (read a valid entry for) the
// given cache key.
func touchKey(key string) {
	touchedKeysMutex.Lock()
	defer touchedKeysMutex.Unlock()

	touchedKeys = append(touchedKeys, key)
}

// FlushTouches persists the last_touched time of every cache entry that this
// build read. It must be called once at the end of the build.
func FlushTouches() {
	touchedKeysMutex.Lock()
	keys := touchedKeys
	touchedKeys = nil
	touchedKeysMutex.Unlock()

	if len(keys) == 0 {
		return
	}

	conn := dbConnection()
	if conn == nil {
		return
	}

	tx, err := conn.Begin()
	if err != nil {
		logger.PrintWarning("failed to open a build cache transaction: " + err.Error())
		return
	}

	committed := false
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()

	statement, err := tx.Prepare(`UPDATE ` + buildCacheTableName + ` SET last_touched = ? WHERE in_out_mod = ?;`)
	if err != nil {
		logger.PrintWarning("failed to prepare the build cache touch update: " + err.Error())
		return
	}

	defer statement.Close()

	for _, key := range keys {
		if _, err := statement.Exec(buildTime(), key); err != nil {
			logger.PrintWarning("failed to touch build cache entry: " + err.Error())
		}
	}

	if err := tx.Commit(); err != nil {
		logger.PrintWarning("failed to commit build cache touches: " + err.Error())
		return
	}

	committed = true
}

// CacheKey returns the cache key for a compilation unit.
//
// The key is a content fingerprint of the page's source and every file that
// the compiled output depends on (layouts, style/script sidecars, and
// transitively imported components and ESM modules).
//
// Previous implementations keyed on the input file's modification time, which
// meant that changes to any dependency (e.g. a shared __layout.html) never
// invalidated the page, and edits that landed inside the same filesystem
// timestamp tick could be served stale. Hashing the content of every compile
// time input fixes both classes of staleness.
func CacheKey(inputPath string, outputPath string) (string, error) {
	dependencies := discoverDependencies(inputPath)
	fingerprintPaths := append([]string{inputPath}, dependencies...)

	hash := sha256.New()
	for _, fingerprintPath := range fingerprintPaths {
		content, err := filesystem.ReadFile(fingerprintPath)
		if err != nil {
			if fingerprintPath == inputPath {
				// We can't fingerprint a page that we can't read (e.g. when
				// the file was deleted between indexing and compilation).
				return "", err
			}

			// A dependency that vanished between discovery and fingerprinting
			// can't contribute to the fingerprint. The next build will not
			// find it either, so this is stable.
			continue
		}

		hash.Write([]byte(fingerprintPath))
		hash.Write([]byte{0})
		hash.Write(content)
		hash.Write([]byte{0})
	}

	// CLI arguments can change the compiled output of a page without changing
	// any input file, so they are folded into the fingerprint.
	hash.Write([]byte(outputFingerprint()))

	return inputPath + "->" + outputPath + "@" + hex.EncodeToString(hash.Sum(nil)), nil
}

// IsCached returns whether the compiler can skip recompiling the given input
// path because the output already exists and every compile time input (see
// CacheKey) is unchanged since the last build.
func IsCached(inputPath string, outputPath string, key string) bool {
	if key == "" {
		return false
	}

	if _, err := os.Stat(outputPath); err != nil {
		return false
	}

	// Opening the connection loads (or refreshes) the in memory key snapshot.
	// A nil connection means the cache is unavailable (everything misses).
	if dbConnection() == nil {
		return false
	}

	if !hasCachedKey(key) {
		return false
	}

	// A cache hit is a use of the entry: record it (batched, see
	// FlushTouches) so that the vacuum keeps entries that builds still rely
	// on.
	touchKey(key)

	return true
}

// CacheRecord describes one successfully compiled asset to record in the
// build cache.
type CacheRecord struct {
	InputPath  string
	OutputPath string
	Key        string
}

// CacheAssets records successful compilations in the build cache.
//
// All records are written in a single transaction. During a full build one
// record is produced per page, and committing them individually would mean
// one write transaction (and one file lock round trip) per page on a database
// that the parallel compilation workers are also reading from.
//
// The key must be the same key that was checked (and that the page was
// compiled against), so that a dependency changing mid-build can never record
// a cache entry for output that was compiled from the old dependency content.
//
// Cache entries for superseded content of the same input/output pair are
// pruned so that the cache database doesn't grow with every content change.
func CacheAssets(records []CacheRecord) {
	if len(records) == 0 {
		return
	}

	conn := dbConnection()
	if conn == nil {
		return
	}

	tx, err := conn.Begin()
	if err != nil {
		logger.PrintWarning("failed to open a build cache transaction: " + err.Error())
		return
	}

	committed := false
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()

	for _, record := range records {
		// All keys for this input/output pair share this prefix, so pruning on
		// the prefix removes entries for old content versions without touching
		// entries for other pages.
		keyPrefix := record.InputPath + "->" + record.OutputPath + "@"

		// substr() is used instead of LIKE so that special characters in file
		// paths (e.g. "%") can't broaden the match.
		if _, err := tx.Exec(
			`DELETE FROM `+buildCacheTableName+` WHERE substr(in_out_mod, 1, ?) = ? AND in_out_mod != ?;`,
			len(keyPrefix), keyPrefix, record.Key,
		); err != nil {
			logger.PrintWarning("failed to prune old build cache entries for " + record.InputPath)
			continue
		}

		// The delete above only prunes keys for *other* content versions of
		// this input/output pair. When the same key is recorded again (e.g.
		// the output directory was deleted but the cache entry is still
		// valid), the insert upserts so that the existing row's touch time is
		// refreshed instead of failing the primary key constraint.
		if _, err := tx.Exec(
			`INSERT INTO `+buildCacheTableName+` (in_out_mod, last_touched) VALUES (?, ?)
				ON CONFLICT(in_out_mod) DO UPDATE SET last_touched = excluded.last_touched;`,
			record.Key, buildTime(),
		); err != nil {
			logger.PrintWarning("failed to record build cache entry for " + record.InputPath + ": " + err.Error())
		}
	}

	if err := tx.Commit(); err != nil {
		logger.PrintWarning("failed to commit build cache entries: " + err.Error())
		return
	}

	committed = true

	// Keep the in memory snapshot in sync with the committed entries so that
	// (e.g. in listen mode) a later compilation of the same build sees the
	// freshly recorded entries.
	for _, record := range records {
		rememberCachedKey(record.Key)
	}
}

// outputFingerprint captures the CLI arguments that change a page's compiled
// output without changing any input file.
func outputFingerprint() string {
	args := cli.GetArgs()
	return fmt.Sprintf(
		"format=%v;devtools=%v;no-reactivity=%v;isolated-pages=%v",
		args.WithFormatting,
		args.HasDevTools,
		args.NoReactivity,
		args.IsolatedPages,
	)
}
