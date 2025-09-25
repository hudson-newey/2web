package rust

import (
	"hudson-newey/2web/src/content/css"
)

type rustCode = string

type RustFile struct {
	Content rustCode
}

func (model *RustFile) ToWasm(filePath string) css.CSSFile {
	content := model.wasmContent(filePath)
	return css.CSSFile{Content: content}
}

func (model *RustFile) wasmContent(filePath string) string {
	return ""
}
