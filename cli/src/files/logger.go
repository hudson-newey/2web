package files

import "fmt"

func logCreation(model File) {
	fmt.Printf("\033[36mCreate\033[0m (\033[32mSuccess\033[0m):\t%s\n", model.Path)
}

func logMigration(model Migration) {
	fmt.Printf("\033[34mMigrate\033[0m (\033[32mSuccess\033[0m):\t%s\n", model.TargetPath)
}
