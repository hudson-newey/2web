package files

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/hudson-newey/2web-cli/src/cli"
)

const directoryPerms os.FileMode = os.ModePerm
const filePerms os.FileMode = 0644

func WriteFiles(files []File) {
	for _, file := range files {
		recursiveWriteFile(file)
	}
}

func recursiveWriteFile(model File) {
	if model.IsDirectory {
		createDirectory(model)
	} else {
		createFile(model)
	}

	logCreation(model)
}

func createDirectory(dirModel File) {
	err := os.Mkdir(dirModel.Path, directoryPerms)
	if err != nil {
		log.Fatal(err)
	}

	if len(dirModel.Children) > 0 {
		WriteFiles(dirModel.Children)
	}
}

func createFile(fileModel File) {
	// We check if CopyFromPath is set instead of Content because it is possible
	// to create a file with empty content, however, it's not possible to copy
	// from an empty path.
	if fileModel.CopyFromPath != "" {
		pathSrc, err := exec.LookPath(fileModel.CopyFromPath)
		if err != nil {
			errorMsg := fmt.Sprintf("Could not find binary in PATH: '%s'", fileModel.CopyFromPath)
			cli.PrintWarning(errorMsg)
		}

		CopyPath(pathSrc, fileModel.Path)
	} else {
		err := os.WriteFile(fileModel.Path, []byte(fileModel.Content), filePerms)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	}
}
