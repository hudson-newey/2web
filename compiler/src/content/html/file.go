package html

import (
	"github.com/yosssi/gohtml"
)

type htmlCode = string

func NewHtmlFile(content htmlCode) *HTMLFile {
	return &HTMLFile{Content: content}
}

type HTMLFile struct {
	Content htmlCode
}

func (model *HTMLFile) AddContent(contentPartial htmlCode) {
	model.Content += contentPartial
}

func (model *HTMLFile) Format() {
	model.Content = gohtml.Format(model.Content)
}
