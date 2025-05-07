package templates

import (
	"fmt"
	"os"

	"github.com/hudson-newey/2web-cli/src/files"
)

func AspectTemplate(name string) {
	dirPath := fmt.Sprintf("src/aspects/%s/", name)

	// ignore errors from this because we expect this to fail (because the
	// directory) already exists
	os.Mkdir("src/aspects/", os.ModePerm)

	templateFiles := []files.File{
		{
			Path:        dirPath,
			IsDirectory: true,
			Children: []files.File{
				{
					Path:    dirPath + name + ".aspect.js",
					Content: createAspectContent(name),
				},
				{
					Path:    dirPath + name + ".aspect.spec.js",
					Content: createAspectTestContent(name),
				},
			},
		},
	}

	files.WriteFiles(templateFiles)
}

func createAspectContent(name string) string {
	camelAspectName := capitalizeFirst(name)

	return fmt.Sprintf(`class %s {
}
`, camelAspectName)
}

func createAspectTestContent(name string) string {
	return fmt.Sprintf(`describe("%sAspect", () => {
});
`, name)
}
