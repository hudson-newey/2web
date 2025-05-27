package templates

import (
	"fmt"
	"os"

	"github.com/hudson-newey/2web-cli/src/files"
)

func GuardTemplate(guardName string) {
	guardPath := fmt.Sprintf("src/guards/%s/", guardName)

	// ignore errors from this because we expect this to fail (because the
	// directory) already exists
	os.Mkdir("src/components/", os.ModePerm)

	templateFiles := []files.File{
		{
			Path:        guardPath,
			IsDirectory: true,
			Children: []files.File{
				{
					Path:        guardPath + guardName + ".guard.ts",
					Content:     createGuardContent(guardName),
					IsDirectory: false,
				},
				{
					Path:        guardPath + guardName + ".guard.spec.ts",
					Content:     createGuardTestContent(guardName),
					IsDirectory: false,
				},
			},
		},
	}

	files.WriteFiles(templateFiles)
}

func createGuardContent(name string) string {
	return fmt.Sprintf(`
	import { RouteGuard } from "@two-web/kit/route-guards";

	function failCondition(): boolean {
		return false;
	}

	const %s = () => {
		if (failCondition()) {
		  // redirect to the login page
			return "/login";
		}
	};
`, name)
}

func createGuardTestContent(name string) string {
	return fmt.Sprintf(`describe("%s", () => {
});
`, name)
}
