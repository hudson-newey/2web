package compiler

import (
	"fmt"
	"hudson-newey/2web/src/ssg"
	"log"
	"os"
)

func CompileFile(path string) string {
	log.Println("\t-", path)

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return ""
	}

	ssgSource := string(data)
	stable := false
	for {
		ssgSource, stable = ssg.ProcessStaticSite(ssgSource, path)

		if stable {
			break
		}
	}

	return compileSource(ssgSource)
}

func compileSource(fileContent string) string {
	return fileContent
}
