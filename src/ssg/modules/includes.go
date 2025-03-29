package modules

import (
	"os"
	"strings"
)

func IncludeSsgContent(value string, filePath string) string {
	hostDirectoryEnd := strings.LastIndex(filePath, "/")
	hostDirectory := filePath[:hostDirectoryEnd]

	data, err := os.ReadFile(hostDirectory + value)
	if err != nil {
		panic(err)
	}

	return string(data)
}
