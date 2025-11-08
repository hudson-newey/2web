package twoscript

import "os"

type twoScriptCode = string

func NewTwoScriptFile() *TwoScriptFile {
	return &TwoScriptFile{}
}

func FromFilePath(filePath string) *TwoScriptFile {
	content, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	newTwoScriptFile := NewTwoScriptFile()
	newTwoScriptFile.AddContent(string(content))

	return newTwoScriptFile
}

func FromContent(content twoScriptCode) *TwoScriptFile {
	newTwoScriptFile := NewTwoScriptFile()
	newTwoScriptFile.AddContent(content)

	return newTwoScriptFile
}

// TODO: This is probably why Svelte migrated to runes...
// I should probably refactor the 2web reactive scripts to use "rune like"
// syntax so that I don't have to maintain a separate model for
// parsing/injecting reactive code.
type TwoScriptFile struct {
	Content twoScriptCode
}

func (model *TwoScriptFile) AddContent(partialContent string) {
	model.Content += partialContent
}
