package templates

import (
	"fmt"
	"os"

	"github.com/hudson-newey/2web-cli/src/files"
)

func ModelTemplate(modelName string) {
	servicePath := fmt.Sprintf("src/models/%s/", modelName)

	// ignore errors from this because we expect this to fail (because the
	// directory) already exists
	os.Mkdir("src/models/", os.ModePerm)

	templateFiles := []files.File{
		{
			Path:        servicePath,
			IsDirectory: true,
			Children: []files.File{
				{
					Path:    servicePath + modelName + ".model.ts",
					Content: createModelContent(modelName),
				},
				{
					Path:    servicePath + modelName + ".model.spec.ts",
					Content: createModelTestContent(modelName),
				},
			},
		},
	}

	files.WriteFiles(templateFiles)
}

func createModelContent(name string) string {
	return fmt.Sprintf(`export class %sModel {
	constructor() {
	}
}
`, name)
}

func createModelTestContent(name string) string {
	return fmt.Sprintf(`describe("%sModel", () => {
});
`, name)
}
