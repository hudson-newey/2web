package page

import "hudson-newey/2web/src/utils"

func (model *Page) Write(filePath string) {
	utils.WriteFile(model.Html.Content, filePath)

	for _, file := range model.Css {
		utils.WriteFile(file.RawContent(), file.OutputPath())
	}

	for _, file := range model.JavaScript {
		if file.IsCompilerOnly() {
			continue
		}

		utils.WriteFile(file.RawContent(), file.OutputPath())
	}
}
