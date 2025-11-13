package docx

import (
	"os"
)

func NewDocxFile(inputPath string) DocxFile {
	// We use os.ReadFile instead of filesystem.ReadFile here because a docx file
	// is likely to be an entry point that is only read once during a build
	// because is a content page that won't be reused elsewhere.
	// By using os.ReadFile, we avoid the memory overhead of caching the file
	// and the processing overhead of mutex locks.
	content, err := os.ReadFile(inputPath)
	if err != nil {
		return DocxFile{}
	}

	return DocxFile{
		InputPath: inputPath,
		Data:      content,
	}
}

type DocxFile struct {
	InputPath string
	Data      []byte
}
