package html

type htmlCode = string

type HTMLFile struct {
	content htmlCode
}

func (model *HTMLFile) AddContent(contentPartial htmlCode) {
	model.content += contentPartial
}
