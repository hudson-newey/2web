package css

type cssCode = string

type CSSFile struct {
	content cssCode
}

func (model *CSSFile) AddContent(contentPartial cssCode) {
	model.content += contentPartial
}
