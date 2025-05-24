package html

type htmlCode = string

type HTMLFile struct {
	Content htmlCode
}

func (model *HTMLFile) AddContent(contentPartial htmlCode) {
	model.Content += contentPartial
}
