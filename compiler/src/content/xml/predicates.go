package xml

import "strings"

func IsXmlFile(filePath string) bool {
	return strings.HasSuffix(filePath, ".xml")
}
