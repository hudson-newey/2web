package generators

import (
	"fmt"
	"os"

	"github.com/hudson-newey/2web-cli/src/files"
)

func InterceptorGenerator(name string) {
	dirPath := fmt.Sprintf("src/interceptors/%s/", name)

	// ignore errors from this because we expect this to fail (because the
	// directory) already exists
	os.Mkdir("src/interceptors/", os.ModePerm)

	templateFiles := []files.File{
		{
			Path:        dirPath,
			IsDirectory: true,
			Children: []files.File{
				{
					Path:    dirPath + name + ".interceptor.ts",
					Content: createInterceptorContent(name),
				},
				{
					Path:    dirPath + name + ".interceptor.spec.ts",
					Content: createInterceptorTestContent(name),
				},
			},
		},
	}

	files.WriteFiles(templateFiles)
}

// TODO: I should improve this interceptor code. It should not be monkey
// patching the window.fetch function
func createInterceptorContent(name string) string {
	return fmt.Sprintf(`function %sInterceptor(...args) {
	return args
}

const originalFetch = window.fetch;
function patchedFetch(...args) {
	const transformedReq = transform(...args);
	return originalFetch(...transformedReq);
}

window.fetch = patchedFetch;
`, name)
}

func createInterceptorTestContent(name string) string {
	return fmt.Sprintf(`describe("%sInterceptor", () => {
});
`, name)
}
