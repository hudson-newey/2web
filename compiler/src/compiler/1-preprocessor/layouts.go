package preprocessor

import (
	"os"
	"path/filepath"
	"strings"
)

var cachedLayouts = map[string]string{}

func expandLayout(filePath string, content *string) {
	layoutDir := filepath.Dir(filePath)
	layoutFile := filepath.Join(layoutDir, "__layout.html")

	layoutString := ""
	if _, exists := cachedLayouts[layoutFile]; exists {
		layoutString = cachedLayouts[layoutFile]
	} else {
		layoutBytes, err := os.ReadFile(layoutFile)
		if err != nil {
			return
		}

		layoutString = string(layoutBytes)
		cachedLayouts[layoutFile] = layoutString
	}

	// We work with a lot of byte slices here so that we don't need to convert the
	// layout to a string and back again multiple times.
	slotSelector := "<slot></slot>"
	slotIndex := strings.Index(layoutString, slotSelector)
	if slotIndex < 0 {
		return
	}

	replaced := strings.Replace(
		layoutString,
		slotSelector,
		*content,
		1,
	)

	*content = replaced
}
