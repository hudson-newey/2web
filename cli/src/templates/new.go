package templates

import (
	"fmt"

	"github.com/hudson-newey/2web-cli/src/files"
)

const twoWebVersion = "latest"

const indexHtmlContent = `<!DOCTYPE html>
<title>My New Website!</title>
<link rel="stylesheet" href="styles/home.css" />

<script compiled>
$ count = 0;
</script>

<h1>My New 2Web Website</h1>

<span>{{ $count }}</span>

<button @click="$count = $count + 1">+1</button>
<button @click="$count = $count - 1">-1</button>
`

var packageJsonContent string = fmt.Sprintf(`{
  "name": "2web-example-project",
  "version": "0.1.0",
  "description": "Add your description",
  "type": "module",
  "scripts": {
    "dev": "2web serve",
    "build": "2web build",
    "lint": "2web lint"
  },
  "author": "your-name",
  "license": "your-license",
  "dependencies": {
    "@two-web/kit": "%s"
  },
  "devDependencies": {
    "@two-web/compiler": "%s",
    "@two-web/cli": "%s",
    "@web/test-runner": "^0.20.2",
    "vite": "^6.3.5",
    "typescript": "^5.8.3",
    "@biomejs/biome": "^2.2.2"
  }
}
`, twoWebVersion, twoWebVersion, twoWebVersion)

const tsconfigContent = `{
  "extends": "@two-web/cli/templates/tsconfig.json",
}
`

const readmeContent = `# 2Web Example Project

## Starting your project

` + "```sh" + `
$ 2web serve
>
` + "```" + `

Since this environment uses Vite for the development server, reload speeds
should be extremely fast and have support for  hot module replacement (HMR).

## Building your project

` + "```sh" + `
$ 2web build
>
` + "```" + `

## Testing your project

` + "```" + `
$ 2web test
>
` + "```" + `
`

const gitignoreContent = `dist/
node_modules/
tmp/
.vite/

vite.config.ts.timestamp-*.mjs

.env
.env.development.local
.env.test.local
.env.production.local
.env.local
`

const vscodeSettingsContent = `{
	"files.eol": "\n",
	"files.insertFinalNewline": true,
	"files.trimTrailingWhitespace": true,
	"editor.tabSize": 2,
	"editor.defaultFormatter": "biomejs.biome",
	["2web"]: {
		"files.autoSave": "afterDelay",
		"files.autoSaveDelay": 0,
	},
}`

const viteSettingsContent = `import { defineConfig } from 'vite'
import twoWeb from "@two-web/kit/vite-plugin";

export default defineConfig({
  appType: "mpa",
	plugins: [twoWeb()],
})
`

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
					Path:        path + "/tsconfig.json",
					Content:     tsconfigContent,
					IsDirectory: false,
				},
				{
					Path:        path + "/vite.config.ts",
					Content:     viteSettingsContent,
					IsDirectory: false,
				},
				{
					Path:        path + "/.vscode/",
					IsDirectory: true,
					Children: []files.File{
						{
							Path:        path + "/.vscode/settings.json",
							Content:     vscodeSettingsContent,
							IsDirectory: false,
						},
					},
				},
				{
					Path:        path + "/src/",
					IsDirectory: true,
					Children: []files.File{
						{
							Path:        path + "/src/index.html",
							Content:     indexHtmlContent,
							IsDirectory: false,
						},
						{
							Path:        path + "/src/styles/",
							IsDirectory: true,
							Children: []files.File{
								{
									Path:        path + "/src/styles/home.css",
									Content:     "",
									IsDirectory: false,
								},
							},
						},
					},
				},
			},
		},
	}

	files.WriteFiles(templateFiles)
}
