package filesystem

import (
	"os"
	"path/filepath"
)

func CreateFile(outputPath string) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	return file.Close()
}
