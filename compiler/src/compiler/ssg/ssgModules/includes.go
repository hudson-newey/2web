package ssgModules

import (
	"fmt"
	"hudson-newey/2web/src/content/document/documentErrors"
	"hudson-newey/2web/src/models"
	"os"
	"strings"
)

func IncludeSsgContent(value string, filePath string) string {
	hostDirectoryEnd := strings.LastIndex(filePath, "/")
	hostDirectory := filePath[:hostDirectoryEnd]

	data, err := os.ReadFile(hostDirectory + value)
	if err != nil {
		documentErrors.AddError(models.Error{
			FilePath: filePath,
			Message:  fmt.Sprintf("Failed to include file: %s\n%s", value, err.Error()),
		})
		return ""
	}

	return string(data)
}
