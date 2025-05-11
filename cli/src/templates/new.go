package templates

import "github.com/hudson-newey/2web-cli/src/files"

const indexHtmlContent = `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>My New Website!</title>
    <link rel="stylesheet" href="styles/home.css" />
  </head>

  <body>
    <script compiled>
    $ count = 0;
    </script>

    <main>
      <h1>My New 2Web Website</h1>

      <span>{{ $count }}</span>

      <button @click="$count = $count + 1">+1</button>
      <button @click="$count = $count - 1">-1</button>
    </main>
  </body>
</html>
`

const packageJsonContent = `{
  "name": "2web-example-project",
  "version": "0.1.0",
  "description": "Add your description",
  "type": "module",
  "scripts": {
    "dev": "vite ./dist/"
  },
  "author": "your-name",
  "license": "your-license",
  "devDependencies": {
    "vite": "^6.2.6"
  }
}
`

const viteConfigContent = `
import { defineConfig } from "vite";

export default defineConfig({
  appType: "mpa",
  base: "",
});
`

const readmeContent = `# 2Web Example Project

## Setting up your development environment

You can build your 2web project by simply running the 2web compiler over the
src/ directory.

An example can be seen below.

` + "```sh" + `
$ 2web -i src/ -o dist/
>
` + "```" + `

Since this environment uses Vite for the development server, reload speeds
should be extremely fast and have support for  hot module replacement (HMR).

You can start the vite dev server with the following command.

` + "```sh" + `
$ npm run dev
>
` + "```" + `
`

func NewTemplate(path string) {
	templateFiles := []files.File{
		{
			Path:        path + "/",
			IsDirectory: true,
			Children: []files.File{
				{
					Path:        path + "/README.md",
					Content:     readmeContent,
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
					Path:        path + "/vite.config.js",
					Content:     viteConfigContent,
					IsDirectory: false,
				},
				{
					Path:        path + "/src/",
					IsDirectory: true,
					Children: []files.File{
						{
							Path:    path + "/src/index.html",
							Content: indexHtmlContent,
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
