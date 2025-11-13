package pdf

import (
	"fmt"
	"hudson-newey/2web/src/cli"
	"os"
)

func FromContent(inputPath string, data []byte) *PdfFile {
	return &PdfFile{
		InputPath: inputPath,
		Content:   data,
	}
}

func FromFilePath(inputPath string) *PdfFile {
	return &PdfFile{
		InputPath: inputPath,
	}
}

type PdfFile struct {
	InputPath string
	Content   []byte
}

func (model *PdfFile) Data() []byte {
	if model.Content != nil {
		return model.Content
	}

	// We use os.ReadFile instead of filesystem.ReadFile here because a pdf file
	// is likely to be an entry point that is only read once during a build
	// because is a content page that won't be reused elsewhere.
	// By using os.ReadFile, we avoid the memory overhead of caching the file
	// and the processing overhead of mutex locks.
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
