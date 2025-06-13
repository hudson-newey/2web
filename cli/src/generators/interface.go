package generators

import (
	"fmt"

	"github.com/hudson-newey/2web-cli/src/files"
)

func InterfaceGenerator(interfaceName string) {
	// Interfaces are typically scoped to a specific model.
	// In this case, we typically want them to be sidecar'd next to the model that
	// they are associated with.
	// Therefore, unlike other generators, an interface does not have its own
	// directory.
	interfacePath := fmt.Sprintf("src/%s", interfaceName)

	templateFiles := []files.File{
		{
			Path:        interfacePath,
			Content:     createInterfaceContent(interfaceName),
			IsDirectory: false,
		},
	}

	files.WriteFiles(templateFiles)
}

func createInterfaceContent(name string) string {
	return fmt.Sprintf(`export interface %s {
}`, name)
}
