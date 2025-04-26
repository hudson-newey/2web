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

      <span *innerText="$count"></span>

      <button @click="$count = $count + 1">+1</button>
      <button @click="$count = $count - 1">-1</button>
    </main>
  </body>
</html>
`

func NewTemplate(path string) {
	templateFiles := []files.File{
		{
			Path:        path + "/",
			IsDirectory: true,
			Children: []files.File{
				{
					Path:        path + "/README.md",
					Content:     path + "\n",
					IsDirectory: false,
				},
				{
					Path:        path + "/LICENSE",
					Content:     "",
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
