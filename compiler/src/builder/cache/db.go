package cache

import (
	"database/sql"
	"hudson-newey/2web/src/cli"
	"hudson-newey/2web/src/filesystem"
	"os"
	"path"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const buildCacheTableName = "build_records"

// The build thread pool compiles pages in parallel, so the shared connection
// and every read of it must be guarded by a mutex.
var dbMutex sync.Mutex

var cachedConnection *sql.DB = nil

// This is public so that the main function can close the connection at the end
// of the programs execution.
func CloseDBConnection() {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	if cachedConnection != nil {
		cachedConnection.Close()
		cachedConnection = nil
	}
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

	db.SetMaxOpenConns(8)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	// The table creation runs on the already locked connection. Calling
	// dbConnection() here would deadlock on dbMutex.
	createDbTables(db)

	return cachedConnection
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
