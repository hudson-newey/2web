package scss

import (
	"fmt"
	"hudson-newey/2web/src/content/css"

	sass "github.com/bep/godartsass/v2"
)

type sassCode = string

type SassFile struct {
	Content sassCode
}

func (model *SassFile) ToCss(filePath string) css.CSSFile {
	content := model.cssContent(filePath)
	return css.CSSFile{Content: content}
}

func (model *SassFile) cssContent(filePath string) string {
	transpiler, err := sass.Start(sass.Options{})
	if err != nil {
		fmt.Println("dart sass may not be installed on your system. Try npm install -g sass")
		panic(err)
	}

	cssContent, err := transpiler.Execute(sass.Args{
		Source:       model.Content,
		IncludePaths: []string{filePath},
	})

	if err != nil {
		panic(err)
	}

	return cssContent.CSS
}
