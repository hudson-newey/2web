package models

import (
	"errors"
	"hudson-newey/2web/src/compiler/lexer"
	"os"
	"strings"
)

type Component struct {
	// The selector that can be used in the template to reference this component
	// e.g. "Footer" for <Footer />
	DomSelector string
	ImportPath  string
	Node        *lexer.LexNode[lexer.ImportNode]
}

func (model *Component) HtmlContent(workingPath string) (string, error) {
	hostDirectoryEnd := strings.LastIndex(workingPath, "/")
	hostDirectory := workingPath[:hostDirectoryEnd]

	componentPath := hostDirectory + model.ImportPath

	data, err := os.ReadFile(componentPath)
	if err != nil {
		return "", errors.New("could not resolve import " + componentPath)
	}

	return string(data), nil
}
