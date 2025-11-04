package pdf

import (
	"fmt"
	"hudson-newey/2web/src/cli"
	"os"
)

func NewPdfFile(inputPath string) *PdfFile {
	return &PdfFile{
		InputPath: inputPath,
	}
}

type PdfFile struct {
	InputPath string
}

func (model *PdfFile) Data() []byte {
	data, err := os.ReadFile(model.InputPath)
	if err != nil {
		return []byte{}
	}

	return data
}

func (model *PdfFile) OutputPath() string {
	inPath := *cli.GetArgs().InputPath
	outPath := *cli.GetArgs().OutputPath

	// Remove the input path from the beginning of the file path to get the
	// relative path.
	//
	// TODO: I think this breaks down in single file mode.
	basename := model.InputPath[len(inPath):]

	return fmt.Sprintf("%s/%s", outPath, basename)
}
