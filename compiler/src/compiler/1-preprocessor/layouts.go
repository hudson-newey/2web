package preprocessor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func expandLayout(filePath string, content string) string {
	layoutDirectory := filepath.Dir(filePath)
	layoutFile := fmt.Sprintf("%s/__layout.html", layoutDirectory)

	if _, err := os.Stat(layoutFile); err == nil {
		layoutContentBytes, err := os.ReadFile(layoutFile)
		if err != nil {
			panic(err)
		}

		layoutContent := string(layoutContentBytes)

		// Replace the <slot></slot> tag in the layout with the content.
		expandedContent := strings.ReplaceAll(layoutContent, "<slot></slot>", content)
		return expandedContent
	}

	return content
}
