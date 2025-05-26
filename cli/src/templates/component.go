package templates

import (
	"fmt"
	"os"

	"github.com/hudson-newey/2web-cli/src/files"
)

func ComponentTemplate(componentName string) {
	componentPath := fmt.Sprintf("src/components/%s/", componentName)

	// ignore errors from this because we expect this to fail (because the
	// directory) already exists
	os.Mkdir("src/components/", os.ModePerm)

	templateFiles := []files.File{
		{
			Path:        componentPath,
			IsDirectory: true,
			Children: []files.File{
				{
					Path:        componentPath + componentName + ".component.html",
					Content:     createComponentContent(componentName),
					IsDirectory: false,
				},
				{
					Path:        componentPath + componentName + ".component.spec.ts",
					Content:     createComponentTestContent(componentName),
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

func createComponentTestContent(name string) string {
	return fmt.Sprintf(`describe("%s", () => {
});
`, name)
}
