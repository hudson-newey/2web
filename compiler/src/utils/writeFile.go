package utils

import (
	"os"
	"path/filepath"
)

// Overwrites or creates a file at the given path with the given content and
// creates any necessary directories.
func WriteFile(content string, outputPath string) {
	if err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm); err != nil {
		panic(err)
	}

	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		panic(err)
	}
}

func WriteBinaryFile(content []byte, outputPath string) {
	if err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm); err != nil {
		panic(err)
	}

	if err := os.WriteFile(outputPath, content, 0644); err != nil {
		panic(err)
	}
}

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
