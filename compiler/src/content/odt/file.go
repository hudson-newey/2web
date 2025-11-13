package odt

import "os"

func NewOdtFile(inputPath string) OdtFile {
	// We use os.ReadFile instead of filesystem.ReadFile here because an odt file
	// is likely to be an entry point that is only read once during a build
	// because is a content page that won't be reused elsewhere.
	// By using os.ReadFile, we avoid the memory overhead of caching the file
	// and the processing overhead of mutex locks.
	content, err := os.ReadFile(inputPath)
	if err != nil {
		return OdtFile{}
	}

	return OdtFile{
		InputPath: inputPath,
		Data:      content,
	}
}

type OdtFile struct {
	InputPath string
	Data      []byte
}
