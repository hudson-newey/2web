package odt

import "os"

func NewOdtFile(inputPath string) OdtFile {
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
