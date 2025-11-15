package templates

import (
	_ "embed"
	"fmt"

	"github.com/hudson-newey/2web-cli/src/files"
)

// We use a lot of embedded strings here to represent the template files.
// Because the "new" command is typically run before any files exist, we can't
// use the `two-web/cli/templates` package that we would normally use for
// templates.
// I have included them in separate files for easier maintenance and editor /
// linting support.

const twoWebVersion = "latest"

//go:embed newStatic/index.html
var indexHtmlContent string

//go:embed newStatic/index.spec.ts
var indexTestContent string

//go:embed newStatic/__style.css
var routeStylesContent string

//go:embed newStatic/components/counter/counter.component.html
var counterComponentContent string

// We cannot use go embed here because we need to inject the version.
var packageJsonContent string = fmt.Sprintf(`{
  "name": "2web-example-project",
  "version": "0.1.0",
  "description": "Add your description",
  "type": "module",
  "scripts": {
    "dev": "2web serve",
    "build": "2web build",
    "lint": "2web lint",
    "format": "2web format",
    "test": "2web test"
  },
  "author": "your-name",
  "license": "your-license",
  "dependencies": {
    "@two-web/kit": "%s"
  },
  "devDependencies": {
    "@two-web/compiler": "%s",
    "@two-web/cli": "%s",
    "playwright": "^1.56.1",
    "typescript": "^5.9.3",
    "oxlint": "^1.28.0",
    "prettier": "^3.6.2"
  }
}
`, twoWebVersion, twoWebVersion, twoWebVersion)

// The tsconfig and gitignore files have a weird name because I don't want them
// interfering with this projects own tsconfig or gitignore.

//go:embed newStatic/tsconfig
var tsconfigContent string

//go:embed newStatic/gitignore
var gitignoreContent string

//go:embed newStatic/README.md
var readmeContent string

//go:embed newStatic/playwright.config.ts
var playwrightConfigContent string

//go:embed newStatic/.vscode/settings.json
var vscodeSettingsContent string

//go:embed newStatic/.vscode/launch.json
var vscodeLaunchContent string

//go:embed newStatic/.vscode/tasks.json
var vscodeTasksContent string

//go:embed newStatic/.vscode/extensions.json
var vscodeExtensionsContent string

//go:embed newStatic/.vscode/mcp.json
var vscodeMcpContent string

//go:embed newStatic/.github/copilot-instructions.md
var copilotInstructions string

func NewTemplate(path string) {
	templateFiles := []files.File{
		{
			Path:        path + "/",
			IsDirectory: true,
			Children: []files.File{
				{
					Path:        path + "/.gitignore",
					Content:     gitignoreContent,
					IsDirectory: false,
				},
				{
					Path:        path + "/LICENSE",
					Content:     "",
					IsDirectory: false,
				},
				{
					Path:        path + "/package.json",
					Content:     packageJsonContent,
					IsDirectory: false,
				},
				{
					Path:        path + "/README.md",
					Content:     readmeContent,
					IsDirectory: false,
				},
				{
					Path:        path + "/playwright.config.ts",
					Content:     playwrightConfigContent,
					IsDirectory: false,
				},
				{
					Path:        path + "/tsconfig.json",
					Content:     tsconfigContent,
					IsDirectory: false,
				},
				{
					Path:        path + "/.github/",
					IsDirectory: true,
					Children: []files.File{
						{
							Path:        path + "/.github/copilot-instructions.md",
							Content:     copilotInstructions,
							IsDirectory: false,
						},
					},
				},
				{
					Path:        path + "/.vscode/",
					IsDirectory: true,
					Children: []files.File{
						{
							Path:        path + "/.vscode/extensions.json",
							Content:     vscodeExtensionsContent,
							IsDirectory: false,
						},
						{
							Path:        path + "/.vscode/launch.json",
							Content:     vscodeLaunchContent,
							IsDirectory: false,
						},
						{
							Path:        path + "/.vscode/mcp.json",
							Content:     vscodeMcpContent,
							IsDirectory: false,
						},
						{
							Path:        path + "/.vscode/settings.json",
							Content:     vscodeSettingsContent,
							IsDirectory: false,
						},
						{
							Path:        path + "/.vscode/tasks.json",
							Content:     vscodeTasksContent,
							IsDirectory: false,
						},
					},
				},
				// By providing the 2web binaries inside the project, we can ensure
				// consistent versions across different environments.
				// This is primarily needed for lower-complexity apps which do not use
				// npm or any package management system.
				// These binaries should probably be removed once node_modules are
				// installed.
				{
					Path:        path + "/bin/",
					IsDirectory: true,
					Children: []files.File{
						{
							Path:         path + "/bin/2web",
							CopyFromPath: "2web",
							IsDirectory:  false,
						},
						{
							Path:         path + "/bin/2webc",
							CopyFromPath: "2webc",
							IsDirectory:  false,
						},
						{
							Path:         path + "/bin/2web-mcp",
							CopyFromPath: "2web-mcp",
							IsDirectory:  false,
						},
					},
				},
				{
					Path:        path + "/src/",
					IsDirectory: true,
					Children: []files.File{
						{
							Path:        path + "/src/components/",
							IsDirectory: true,
							Children: []files.File{
								{
									Path:        path + "/src/components/counter/",
									IsDirectory: true,
									Children: []files.File{
										{
											Path:        path + "/src/components/counter/counter.component.html",
											Content:     counterComponentContent,
											IsDirectory: false,
										},
									},
								},
							},
						},
						{
							Path:        path + "/src/index.html",
							Content:     indexHtmlContent,
							IsDirectory: false,
						},
						{
							Path:        path + "/src/index.spec.ts",
							Content:     indexTestContent,
							IsDirectory: false,
						},
						{
							Path:        path + "/src/__style.css",
							Content:     routeStylesContent,
							IsDirectory: false,
						},
					},
				},
			},
		},
	}

	files.WriteFiles(templateFiles)
}
