package builder

import "path/filepath"

func isLayoutFile(filePath string) bool {
	return filepath.Base(filePath) == "__layout.html"
}
