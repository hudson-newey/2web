package files

import (
	"fmt"
	"log"
	"os"
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
	err := os.WriteFile(fileModel.Path, []byte(fileModel.Content), filePerms)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func logCreation(model File) {
	fmt.Printf("\033[36mCreate\033[0m (\033[32mSuccess\033[0m):\t%s\n", model.Path)
}
