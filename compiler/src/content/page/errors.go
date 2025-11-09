package page

import (
	"hudson-newey/2web/src/content/document/documentErrors"
	"hudson-newey/2web/src/models"
)

type PageErrors struct {
	Errors []*models.Error
}

func (model *PageErrors) AddError(err *models.Error) {
	model.Errors = append(model.Errors, err)
	documentErrors.AddErrors(err)
}

func (model *PageErrors) IsErrorFree() bool {
	return len(model.Errors) == 0
}
