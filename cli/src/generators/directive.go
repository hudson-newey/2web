package generators

import (
	"fmt"
	"os"

	"github.com/hudson-newey/2web-cli/src/files"
)

func DirectiveGenerator(directiveName string) {
	componentPath := fmt.Sprintf("src/directives/%s/", directiveName)

	// ignore errors from this because we expect this to fail (because the
	// directory) already exists
	os.Mkdir("src/directives/", os.ModePerm)

	templateFiles := []files.File{
		{
			Path:        componentPath,
			IsDirectory: true,
			Children: []files.File{
				{
					Path:        componentPath + directiveName + ".directive.ts",
					Content:     createDirectiveContent(directiveName),
					IsDirectory: false,
				},
				{
					Path:        componentPath + directiveName + ".directive.spec.ts",
					Content:     createDirectiveTestContent(directiveName),
					IsDirectory: false,
				},
			},
		},
	}

	files.WriteFiles(templateFiles)
}

func createDirectiveContent(name string) string {
	return fmt.Sprintf(`// 2Web directives are compiler polyfills of the "is" global attribute
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Global_attributes/is
export class %sDirective extends HTMLDivElement {
}

customElements.define("my-%s", %sDirective, { extends: "div" });
`, name, name, name)
}

func createDirectiveTestContent(name string) string {
	return fmt.Sprintf(`describe("%sDirective", () => {
});
`, name)
}
