package preprocessor

import (
	"bytes"
	"os"
	"path/filepath"
)

var cachedLayouts = map[string][]byte{}

func expandLayout(filePath string, content *[]byte) {
	layoutDir := filepath.Dir(filePath)
	layoutFile := filepath.Join(layoutDir, "__layout.html")

	layoutBytes := []byte{}
	if _, exists := cachedLayouts[layoutFile]; exists {
		layoutBytes = cachedLayouts[layoutFile]
	} else {
		layoutBytes, err := os.ReadFile(layoutFile)
		if err != nil {
			return
		}

		cachedLayouts[layoutFile] = layoutBytes
	}

	// We work with a lot of byte slices here so that we don't need to convert the
	// layout to a string and back again multiple times.
	slot := []byte("<slot></slot>")
	slotIndex := bytes.Index(layoutBytes, slot)
	if slotIndex < 0 {
		return
	}

	replaced := bytes.Replace(layoutBytes, slot, *content, 1)

	*content = replaced
}
