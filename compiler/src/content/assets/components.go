package assets

import "strings"

func IsComponent(filePath string) bool {
	return strings.HasSuffix(filePath, ".component.html")
}
