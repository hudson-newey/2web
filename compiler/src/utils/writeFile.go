package utils

import (
	"fmt"
	"hudson-newey/2web/src/cli"
	"os"
	"path/filepath"
)

func WriteFile(content string, outputPath string) {
	if *cli.GetArgs().ToStdout {
		fmt.Println(content)
	} else {
		os.MkdirAll(filepath.Dir(outputPath), os.ModePerm)
		os.WriteFile(outputPath, []byte(content), 0644)
	}
}
