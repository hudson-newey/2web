package cache

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"strings"

	"hudson-newey/2web/src/constants"
)

// The vacuum keeps the build cache from growing without bound.
//
// It only runs when the cache database is larger than maxCacheSize, so normal
// projects (whose cache database holds one small row per compiled asset) never
// pay for it. When it does run, it removes entries in two stages:
//
//  1. Age based: every entry that wasn't touched within cacheRetention is
//     deleted. An entry is "touched" by a build that recorded it (it was
//     compiled) or that read it (it was served from cache), so this only
//     removes entries that builds have been ignoring.
//  2. Size based: if the cache is still over the limit, the oldest entries
//     (by last_touched) are deleted until the cache is back under a low water
//     mark of half the limit.
//
// The deleted space is then reclaimed with VACUUM. Because the deletes are
// rare (only once the cache has actually grown too large), the cost of the
// full rewrite is amortized over many builds.

const (
	// The cache size above which the cache is vacuumed, and the size it is
	// cut back down to.
	defaultMaxCacheSize = int64(256) << 20

	// How long an untouched entry survives once the cache is over the size
	// limit.
	defaultRetentionSeconds = int64(30 * 24 * 3600)
)

var (
	maxCacheSize     = defaultMaxCacheSize
	retentionSeconds = defaultRetentionSeconds

	// vacuumConfigResolved records that the environment overrides have been
	// applied (the environment is constant for the lifetime of the process).
	vacuumConfigResolved = false
)

// resolveVacuumConfig applies the environment overrides for the vacuum policy.
//
// An invalid override is ignored (the defaults are safe), and the resolution
// happens once because the environment is constant for the lifetime of the
// process.
func resolveVacuumConfig() {
	if vacuumConfigResolved {
		return
	}

	vacuumConfigResolved = true

	if override := os.Getenv(constants.EnvCacheMaxSize); override != "" {
		if size, err := strconv.ParseInt(strings.TrimSpace(override), 10, 64); err == nil && size > 0 {
			maxCacheSize = size
		} else {
			log.Printf("ignoring invalid %s value %q", constants.EnvCacheMaxSize, override)
		}
	}

	if override := os.Getenv(constants.EnvCacheMaxAge); override != "" {
		if days, err := strconv.ParseInt(strings.TrimSpace(override), 10, 64); err == nil && days > 0 {
			retentionSeconds = days * 24 * 3600
		} else {
			log.Printf("ignoring invalid %s value %q", constants.EnvCacheMaxAge, override)
		}
	}
}

// Vacuum reclaims space in the build cache when it has grown too large.
//
// It is called once at the end of a build (after the cache writes and touches
// have been flushed) and returns without doing anything when the cache is
// within its size limit.
func Vacuum() {
	resolveVacuumConfig()

	conn := dbConnection()
	if conn == nil {
		return
	}

	// Flush the write ahead log so that the measured size is where the cache
	// actually is.
	checkpoint(conn)

	size := cacheSize()
	if size <= maxCacheSize {
		return
	}

	var rowCount int64
	if err := conn.QueryRow(`SELECT COUNT(*) FROM ` + buildCacheTableName + `;`).Scan(&rowCount); err != nil {
		log.Printf("build cache vacuum failed to count entries: %v", err)
		return
	}

	if rowCount == 0 {
		return
	}

	// The average bytes per row is measured before anything is deleted: once
	// rows are deleted, the freed pages still count towards the file size
	// (until the VACUUM), which would distort the estimate.
	bytesPerRow := size / rowCount
	if bytesPerRow < 1 {
		bytesPerRow = 1
	}

	deleted := int64(0)

	// Stage 1: entries that builds haven't touched within the retention
	// window.
	result, err := conn.Exec(
		`DELETE FROM `+buildCacheTableName+` WHERE last_touched <= ?;`,
		buildTime()-retentionSeconds,
	)
	if err != nil {
		log.Printf("build cache vacuum failed to remove stale entries: %v", err)
		return
	}

	if removed, err := result.RowsAffected(); err == nil {
		deleted += removed
	}

	// Stage 2: cut the cache back down to a low water mark by evicting the
	// oldest entries.
	deleted += evictToLowWater(conn, bytesPerRow)

	if deleted == 0 {
		return
	}

	// Reclaim the space of the deleted rows.
	if _, err := conn.Exec("VACUUM;"); err != nil {
		log.Printf("build cache vacuum failed to reclaim space: %v", err)
	}

	// VACUUM rewrites the database through the write ahead log, so the log has
	// to be flushed back into the database file (and truncated) before the
	// cache size reflects the reclaimed space.
	checkpoint(conn)

	// Drop the deleted keys from the in memory snapshot so that the snapshot
	// doesn't report entries that were just removed. The vacuum runs after
	// the build pool has finished, so this can't race a concurrent lookup.
	dbMutex.Lock()
	cachedKeySnapshot = loadKeySnapshot(conn)
	dbMutex.Unlock()

	log.Printf(
		"vacuumed the build cache: removed %d stale entries (%d MiB -> %d MiB)",
		deleted, size>>20, cacheSize()>>20,
	)
}

// evictToLowWater deletes the oldest entries (by last_touched) until the
// number of remaining rows fits back under the low water mark (half the size
// limit, estimated with the given bytes per row), and returns how many
// entries were deleted.
//
// The row budget is computed once from the pre-deletion measurement: deleting
// rows doesn't shrink the database file until the VACUUM runs, so re-measuring
// the file size during the eviction would only distort the estimate.
func evictToLowWater(conn *sql.DB, bytesPerRow int64) int64 {
	lowWater := maxCacheSize / 2
	targetRows := lowWater / bytesPerRow

	var remaining int64
	if err := conn.QueryRow(`SELECT COUNT(*) FROM ` + buildCacheTableName + `;`).Scan(&remaining); err != nil {
		log.Printf("build cache vacuum failed to count entries: %v", err)
		return 0
	}

	toEvict := remaining - targetRows
	if toEvict <= 0 {
		return 0
	}

	result, err := conn.Exec(
		`DELETE FROM `+buildCacheTableName+` WHERE in_out_mod IN (
			SELECT in_out_mod FROM `+buildCacheTableName+` ORDER BY last_touched ASC, in_out_mod ASC LIMIT ?
		);`,
		toEvict,
	)
	if err != nil {
		log.Printf("build cache vacuum failed to evict old entries: %v", err)
		return 0
	}

	removed, err := result.RowsAffected()
	if err != nil {
		return 0
	}

	return removed
}
