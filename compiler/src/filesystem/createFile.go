package filesystem

import (
	"os"
	"path/filepath"
)

func CreateFile(outputPath string) {
	if err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm); err != nil {
		panic(err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
}
