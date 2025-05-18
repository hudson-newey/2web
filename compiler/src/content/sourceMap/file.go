package sourceMap

type sourceMapContent = string

type SourceMapFile struct {
	content sourceMapContent
}

func (model *SourceMapFile) AddContent(partialContent sourceMapContent) {
	model.content += partialContent
}
