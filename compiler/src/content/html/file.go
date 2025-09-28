package html

import (
	"io"
	"strings"

	"github.com/yosssi/gohtml"
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

func (model *HTMLFile) Format() {
	model.Content = gohtml.Format(model.Content)
}
