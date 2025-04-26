package templates

import (
	"fmt"
	"os"

	"github.com/hudson-newey/2web-cli/src/files"
)

func ComponentTemplate(componentName string) {
	componentsPath := fmt.Sprintf("src/components/%s/", componentName)

	// ignore errors from this because we expect this to fail (because the
	// directory) already exists
	os.Mkdir("src/components/", os.ModePerm)

	templateFiles := []files.File{
		{
			Path:        componentsPath,
			IsDirectory: true,
			Children: []files.File{
				{
					Path:        componentsPath + componentName + ".component.html",
					Content:     createComponentContent(componentName),
					IsDirectory: false,
				},
				{
					Path:        componentsPath + componentName + ".component.spec.js",
					Content:     createTestContent(componentName),
					IsDirectory: false,
				},
			},
		},
	}

	files.WriteFiles(templateFiles)
}

func createComponentContent(name string) string {
	return fmt.Sprintf(`<script compiled>
$ message = "%s works!";
</script>

<p>{{ $message }}</p>
`, name)
}

func createTestContent(name string) string {
	return fmt.Sprintf(`describe("%s", () => {
});
`, name)
}
