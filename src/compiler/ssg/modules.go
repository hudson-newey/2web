package ssg

import (
	"log"
	"os"
	"strings"
)

func includeSsgContent(value string, filePath string) string {
	log.Println("\t\t-", value)

	hostDirectoryEnd := strings.LastIndex(filePath, "/")
	hostDirectory := filePath[:hostDirectoryEnd]

	data, err := os.ReadFile(hostDirectory + value)
	if err != nil {
		panic(err)
	}

	return string(data)
}
