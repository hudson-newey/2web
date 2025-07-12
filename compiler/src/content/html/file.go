package html

import (
	"io"
	"strings"
)

type htmlCode = string

type HTMLFile struct {
	Content htmlCode
}

func (model *HTMLFile) AddContent(contentPartial htmlCode) {
	model.Content += contentPartial
}

func (model *HTMLFile) Reader() io.Reader {
	return strings.NewReader(model.Content)
}
