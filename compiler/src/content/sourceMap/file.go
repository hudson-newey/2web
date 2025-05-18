package sourceMap

type sourceMapContent = string

type SourceMapFile struct {
	Content sourceMapContent
}

func (model *SourceMapFile) AddContent(partialContent sourceMapContent) {
	model.Content += partialContent
}
