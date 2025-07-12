package models

import (
	"errors"
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"os"
	"path/filepath"
	"strconv"
)

type Component struct {
	// The selector that can be used in the template to reference this component
	// e.g. "Footer" for <Footer />
	DomSelector  string
	Identifier   uint64
	ImportedFrom string
	ImportPath   string
	Node         *lexer.LexNode[lexer.ImportNode]
}

func (model *Component) ComponentPath() (string, error) {
	hostDirectory := filepath.Dir(model.ImportedFrom)
	componentPath := filepath.Join(hostDirectory, model.ImportPath)

	if _, err := os.Stat(componentPath); err == os.ErrNotExist {
		return "", errors.New("could not resolve import " + componentPath)
	}

	return componentPath, nil
}

func (model *Component) DomIdentifier() string {
	return strconv.FormatUint(model.Identifier, 36)
}
