package preprocessor

import (
	"os"
	"path/filepath"
	"strings"
)

var cachedLayouts = map[string]string{}

func expandLayout(filePath string, content string) string {
	layoutDir := filepath.Dir(filePath)
	layoutFile := filepath.Join(layoutDir, "__layout.html")

	layoutContent := ""
	if _, exists := cachedLayouts[layoutFile]; exists {
		layoutContent = cachedLayouts[layoutFile]
	} else {
		layoutBytes, err := os.ReadFile(layoutFile)
		if err != nil {
			return content
		}

		layoutContent = string(layoutBytes)
		cachedLayouts[layoutFile] = layoutContent
	}

	slotSelector := "<slot></slot>"

	return strings.Replace(layoutContent, slotSelector, content, 1)
}
