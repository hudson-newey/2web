package templates

import "github.com/hudson-newey/2web-cli/src/files"

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
			},
		},
	}

	files.WriteFiles(templateFiles)
}
