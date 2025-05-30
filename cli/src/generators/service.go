package generators

import (
	"fmt"
	"os"
	"unicode"

	"github.com/hudson-newey/2web-cli/src/files"
)

func ServiceGenerator(serviceName string) {
	servicePath := fmt.Sprintf("src/services/%s/", serviceName)

	// ignore errors from this because we expect this to fail (because the
	// directory) already exists
	os.Mkdir("src/services/", os.ModePerm)

	templateFiles := []files.File{
		{
			Path:        servicePath,
			IsDirectory: true,
			Children: []files.File{
				{
					Path:    servicePath + serviceName + ".service.ts",
					Content: createServiceContent(serviceName),
				},
				{
					Path:    servicePath + serviceName + ".service.spec.ts",
					Content: createServiceTestContent(serviceName),
				},
			},
		},
	}

	files.WriteFiles(templateFiles)
}

func createServiceContent(name string) string {
	camelServiceName := capitalizeFirst(name)

	return fmt.Sprintf(`export function create%s() {
}

export function get%s() {
}

export function update%s() {
}

export function delete%s() {
}
`, camelServiceName, camelServiceName, camelServiceName, camelServiceName)
}

func createServiceTestContent(name string) string {
	return fmt.Sprintf(`describe("%sService", () => {
});
`, name)
}

// TODO: remove this AI generated code
func capitalizeFirst(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
