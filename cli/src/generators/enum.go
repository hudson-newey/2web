package generators

import (
	"fmt"

	"github.com/hudson-newey/2web-cli/src/files"
)

func EnumGenerator(enumName string) {
	// Enums are typically scoped to a specific component, service, or similar.
	// In this case, we typically want them to be sidecar'd next to the file that
	// they are associated with.
	// Therefore, unlike other generators, enums do not have their own directory.
	enumPath := fmt.Sprintf("src/%s", enumName)

	templateFiles := []files.File{
		{
			Path:        enumPath,
			IsDirectory: false,
			Content:     createEnumContent(enumName),
		},
	}

	files.WriteFiles(templateFiles)
}

func createEnumContent(name string) string {
	return fmt.Sprintf(`export enum %s {
}
`, name)
}
