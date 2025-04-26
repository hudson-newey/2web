package templates

import (
	"fmt"
	"os"

	"github.com/hudson-newey/2web-cli/src/files"
)

const serviceContent = `
`

func ServiceTemplate(serviceName string) {
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
					Path:    servicePath + serviceName + ".service.js",
					Content: serviceContent,
				},
				{
					Path:    servicePath + serviceName + ".service.spec.js",
					Content: createServiceTestContent(serviceName),
				},
			},
		},
	}

	files.WriteFiles(templateFiles)
}

func createServiceTestContent(name string) string {
	return fmt.Sprintf(`describe("%s", () => {
});
`, name)
}
