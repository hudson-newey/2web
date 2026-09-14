package cache

import (
	"database/sql"
	"hudson-newey/2web/src/cli"
	"hudson-newey/2web/src/filesystem"
	"log"
	"os"
	"path"
	"strings"
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

// dbDisabled records that the cache database couldn't be opened so that the
// database isn't retried (and failed) on every cache operation.
var dbDisabled bool

// dbConnection returns the shared cache database connection, or nil when the
// cache is unavailable.
//
// The build cache is disposable state, so a database that can't be opened
// (e.g. because the cache directory is read only) degrades to an empty cache
// instead of failing the build.
func dbConnection() *sql.DB {
	if cli.GetArgs().DisableCache {
		return nil
	}

	dbMutex.Lock()
	defer dbMutex.Unlock()

	if cachedConnection != nil {
		return cachedConnection
	}

	if dbDisabled {
		return nil
	}

	dbPath := dbLocation()
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		if err := filesystem.CreateFile(dbPath); err != nil {
			log.Printf("build cache disabled: failed to create cache database: %v", err)
			dbDisabled = true
			return nil
		}
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Printf("build cache disabled: failed to open cache database: %v", err)
		dbDisabled = true
		return nil
	}

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
	if err := createDbTables(db); err != nil {
		log.Printf("build cache disabled: failed to create cache tables: %v", err)
		db.Close()
		dbDisabled = true
		return nil
	}

	cachedKeySnapshot = loadKeySnapshot(db)

	cachedConnection = db

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

func createDbTables(db *sql.DB) error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS ` + buildCacheTableName + ` (
		in_out_mod TEXT PRIMARY KEY,
		last_touched INTEGER NOT NULL DEFAULT 0
	);`

	if _, err := db.Exec(createTableSQL); err != nil {
		return err
	}

	// Databases created before the vacuum feature have no last_touched column.
	// The migration adds it; the "duplicate column" failure for databases that
	// already have it is expected and ignored. Any other failure leaves the
	// cache working (with entries touched at their insertion time) but the
	// vacuum is disabled for this run.
	if _, err := db.Exec(`ALTER TABLE ` + buildCacheTableName + ` ADD COLUMN last_touched INTEGER NOT NULL DEFAULT 0;`); err != nil && !strings.Contains(strings.ToLower(err.Error()), "duplicate column") {
		return err
	}

	return nil
}

// cacheSize returns the size of the build cache database on disk, including
// its write ahead log.
func cacheSize() int64 {
	total := int64(0)

	for _, suffix := range []string{"", "-wal"} {
		if info, err := os.Stat(dbLocation() + suffix); err == nil {
			total += info.Size()
		}
	}

	return total
}

// checkpoint flushes the write ahead log back into the main database file and
// truncates it, so that cacheSize() measures where the cache actually is (and
// deleted rows actually reclaim space).
func checkpoint(conn *sql.DB) {
	if _, err := conn.Exec("PRAGMA wal_checkpoint(TRUNCATE);"); err != nil {
		log.Printf("failed to checkpoint the build cache: %v", err)
	}
}
