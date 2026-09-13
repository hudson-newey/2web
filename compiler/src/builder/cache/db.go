package cache

import (
	"database/sql"
	"hudson-newey/2web/src/cli"
	"hudson-newey/2web/src/filesystem"
	"log"
	"os"
	"path"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

const buildCacheTableName = "build_records"

// The build thread pool compiles pages in parallel, so the shared connection
// and every read of it must be guarded by a mutex.
var dbMutex sync.Mutex

var cachedConnection *sql.DB = nil

// cachedKeySnapshot holds every cache key currently recorded in the build
// cache database.
//
// The table is tiny (one row per compiled input/output pair), so it is loaded
// once per process and kept in memory. Answering cache lookups from memory
// avoids one SQLite round trip (a cgo call) per page during a build.
var cachedKeySnapshot map[string]bool = nil

// This is public so that the main function can close the connection at the end
// of the programs execution.
func CloseDBConnection() {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	if cachedConnection != nil {
		cachedConnection.Close()
		cachedConnection = nil
	}

	cachedKeySnapshot = nil
}

func dbLocation() string {
	return path.Join(cacheLocation(), "build.db")
}

func dbConnection() *sql.DB {
	if cli.GetArgs().DisableCache {
		panic("Attempted to establish database connection with DisabledCache")
	}

	dbMutex.Lock()
	defer dbMutex.Unlock()

	if cachedConnection != nil {
		return cachedConnection
	}

	dbPath := dbLocation()
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		filesystem.CreateFile(dbPath)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}

	cachedConnection = db

	// SQLite allows exactly one writer at a time, so a pool of connections
	// would only add cross connection file locking (and busy retries) between
	// the parallel page compilation workers. A single connection serializes
	// access at the Go layer, which is faster and deterministic.
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	// Write ahead logging lets the (single) writer and readers proceed without
	// blocking each other, and NORMAL synchronous mode only flushes the write
	// ahead log on checkpoints instead of on every transaction. The build cache
	// is disposable state, so the weaker durability guarantee is acceptable.
	for _, pragma := range []string{
		"PRAGMA journal_mode = WAL;",
		"PRAGMA synchronous = NORMAL;",
	} {
		if _, err := db.Exec(pragma); err != nil {
			log.Printf("failed to apply build cache pragma: %v", err)
		}
	}

	// The table creation runs on the already locked connection. Calling
	// dbConnection() here would deadlock on dbMutex.
	createDbTables(db)

	cachedKeySnapshot = loadKeySnapshot(db)

	return cachedConnection
}

// loadKeySnapshot reads every recorded cache key into memory.
func loadKeySnapshot(db *sql.DB) map[string]bool {
	snapshot := map[string]bool{}

	rows, err := db.Query(`SELECT in_out_mod FROM ` + buildCacheTableName + `;`)
	if err != nil {
		// A failed snapshot load must not break the build; the cache just
		// behaves as empty and entries get recompiled.
		return snapshot
	}

	defer rows.Close()

	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			continue
		}

		snapshot[key] = true
	}

	return snapshot
}

// hasCachedKey reports whether the given cache key is recorded in the build
// cache, using the in memory snapshot.
func hasCachedKey(key string) bool {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	if cachedKeySnapshot == nil {
		return false
	}

	return cachedKeySnapshot[key]
}

// rememberCachedKey adds a cache key to the in memory snapshot.
func rememberCachedKey(key string) {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	if cachedKeySnapshot == nil {
		return
	}

	cachedKeySnapshot[key] = true
}

func createDbTables(db *sql.DB) {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS ` + buildCacheTableName + ` (
		in_out_mod TEXT PRIMARY KEY
	);`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		panic(err)
	}
}
