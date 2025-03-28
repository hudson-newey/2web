package compiler

import (
	"fmt"
	"hudson-newey/2web/src/compiler/ssg"
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

	ssgSource := ssg.ProcessStaticSite(string(data), path)
	return compileSource(ssgSource)
}

func compileSource(fileContent string) string {
	return fileContent
}
