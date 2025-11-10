package routing

import "path/filepath"

func IsLayoutFile(filePath string) bool {
	return filepath.Base(filePath) == "__layout.html"
}
