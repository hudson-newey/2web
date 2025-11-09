package cache

import (
	"os"
	"time"
)

// TODO: Refactor this a lot better because it is still possible to break this.
func IsCached(inputPath string, outputPath string) bool {
	_, err := os.Stat(outputPath)
	if os.IsNotExist(err) {
		return false
	}

	inputFileInfo, err := os.Stat(inputPath)
	if err != nil {
		panic(err)
	}

	modTime := inputFileInfo.ModTime()
	queryKey := newInOutModKey(inputPath, outputPath, modTime)

	conn := dbConnection()

	querySQL := `
	SELECT COUNT(*) FROM ` + buildCacheTableName + ` WHERE in_out_mod = ?;
	`

	var count int
	err = conn.QueryRow(querySQL, queryKey).Scan(&count)
	if err != nil {
		return false
	}

	return count > 0
}

func CacheAsset(inputPath string, outputPath string) {
	inputFileInfo, err := os.Stat(inputPath)
	if err != nil {
		panic(err)
	}

	modTime := inputFileInfo.ModTime()
	cacheKey := newInOutModKey(inputPath, outputPath, modTime)

	conn := dbConnection()

	insertSQL := `
	INSERT INTO ` + buildCacheTableName + `(in_out_mod)
	VALUES (?);`

	_, err = conn.Exec(insertSQL, cacheKey)
	if err != nil {
		panic(err)
	}
}

func newInOutModKey(inputPath string, outputPath string, modTime time.Time) string {
	return inputPath + "->" + outputPath + "@" + modTime.String()
}
