package modules

import (
	"log"
	"os"
	"strings"
)

func IncludeSsgContent(value string, filePath string) string {
	log.Println("\t\t-", value)

	hostDirectoryEnd := strings.LastIndex(filePath, "/")
	hostDirectory := filePath[:hostDirectoryEnd]

	data, err := os.ReadFile(hostDirectory + value)
	if err != nil {
		panic(err)
	}

	return string(data)
}
