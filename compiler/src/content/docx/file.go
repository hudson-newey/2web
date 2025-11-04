package docx

import (
	"os"
)

// TODO: AI generated code - needs rewrite
func NewDocxFile(inputPath string) DocxFile {
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
