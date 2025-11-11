package templates

import (
	"fmt"

	"github.com/hudson-newey/2web-cli/src/files"
)

const twoWebVersion = "latest"

const indexHtmlContent = `<title>My New Website!</title>

<script compiled>
  $ count = 0;
</script>

<h1>Welcome to 2Web</h1>

<span>{{ $count }}</span>

<button @click="$count = $count + 1">+1</button>
<button @click="$count = $count - 1">-1</button>
`

const routeStylesContent = `body {
	margin: 0;
}
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

const gitignoreContent = `# 2Web specific ignores
dist/
.cache/

# Logs
logs
*.log
npm-debug.log*
yarn-debug.log*
yarn-error.log*
lerna-debug.log*

# Diagnostic reports (https://nodejs.org/api/report.html)
report.[0-9]*.[0-9]*.[0-9]*.[0-9]*.json

# Coverage directory used by tools like istanbul
coverage
*.lcov

# Compiled binary addons (https://nodejs.org/api/addons.html)
build/Release

# Dependency directories
node_modules/
jspm_packages/

# TypeScript cache
*.tsbuildinfo

# Optional npm cache directory
.npm

# Optional eslint cache
.eslintcache

# Optional stylelint cache
.stylelintcache

# Optional REPL history
.node_repl_history

# Output of 'npm pack'
*.tgz

# dotenv environment variable files
.env
.env.*
!.env.example

# Stores VSCode versions used for testing VSCode extensions
.vscode-test

# Yarn Integrity file
.yarn-integrity

# yarn v3
.pnp.*
.yarn/*
!.yarn/patches
!.yarn/plugins
!.yarn/releases
!.yarn/sdks
!.yarn/versions

# Vite files
vite.config.js.timestamp-*
vite.config.ts.timestamp-*
.vite/
`

const vscodeSettingsContent = `{
	"files.eol": "\n",
	"files.insertFinalNewline": true,
	"files.trimTrailingWhitespace": true,
	"editor.tabSize": 2,
	"editor.defaultFormatter": "biomejs.biome",
	"editor.codeActionsOnSave": {
		"source.fixAll": "explicit",
		"source.organizeImports": "always",
		"source.addMissingImports.ts": "always",
		"source.removeUnusedImports": "always",
	},
	"search.exclude": {
		"dist/": true,
	},
	// "[html]": {
	// 	"files.autoSave": "afterDelay",
	// 	"files.autoSaveDelay": 0,
	// },
	// "[2web]": {
	// 	"files.autoSave": "afterDelay",
	// 	"files.autoSaveDelay": 0,
	// },
	"recommendations": [
		"streetsidesoftware.code-spell-checker",
		"yzhang.markdown-all-in-one",
		"davidanson.vscode-markdownlint",
		"visualstudioexptteam.vscodeintellicode",
		"aaron-bond.better-comments",

		"ms-vscode.vscode-typescript-next",
		"zignd.html-css-class-completion",
		"biomejs.biome",
	]
}`

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
