package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"

	"hudson-newey/2web/src/cli"
	"hudson-newey/2web/src/filesystem"
	"github.com/hudson-newey/2web/_shared/logger"
)

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

	conn := dbConnection()

	querySQL := `
	SELECT COUNT(*) FROM ` + buildCacheTableName + ` WHERE in_out_mod = ?;
	`

	var count int
	err := conn.QueryRow(querySQL, key).Scan(&count)
	if err != nil {
		return false
	}

	return count > 0
}

// CacheAsset records a successful compilation in the build cache.
//
// The key must be the same key that was checked (and that the page was
// compiled against), so that a dependency changing mid-build can never record
// a cache entry for output that was compiled from the old dependency content.
//
// Cache entries for superseded content of the same input/output pair are
// pruned so that the cache database doesn't grow with every content change.
func CacheAsset(inputPath string, outputPath string, key string) {
	if key == "" {
		return
	}

	conn := dbConnection()

	// All keys for this input/output pair share this prefix, so pruning on the
	// prefix removes entries for old content versions without touching entries
	// for other pages.
	keyPrefix := inputPath + "->" + outputPath + "@"

	// substr() is used instead of LIKE so that special characters in file
	// paths (e.g. "%") can't broaden the match.
	_, err := conn.Exec(
		`DELETE FROM `+buildCacheTableName+` WHERE substr(in_out_mod, 1, ?) = ? AND in_out_mod != ?;`,
		len(keyPrefix), keyPrefix, key,
	)
	if err != nil {
		logger.PrintWarning("failed to prune old build cache entries for " + inputPath)
	}

	insertSQL := `
	INSERT INTO ` + buildCacheTableName + `(in_out_mod)
	VALUES (?);`

	_, err = conn.Exec(insertSQL, key)
	if err != nil {
		panic(err)
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
