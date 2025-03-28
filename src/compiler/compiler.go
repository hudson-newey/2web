package compiler

import (
	"fmt"
	"log"
	"os"
)

func CompileFile(path string) string {
	log.Println(path)

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return ""
	}

	return compileSource(string(data))
}

func compileSource(fileContent string) string {
	return fileContent
}
