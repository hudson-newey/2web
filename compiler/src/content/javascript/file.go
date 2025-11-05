package javascript

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"hudson-newey/2web/src/cli"
	"hudson-newey/2web/src/content/document/documentErrors"
	"hudson-newey/2web/src/models"
	"os"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

type javascriptCode = string

func NewJsFile() JSFile {
	return JSFile{}
}

func FromFilePath(filePath string) JSFile {
	content, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	newJsFile := NewJsFile()
	newJsFile.AddContent(string(content))

	return newJsFile
}

type JSFile struct {
	Content          javascriptCode
	memoisedFileName string
}

func (model *JSFile) RawContent() string {
	result := model.Content

	result = strings.ReplaceAll(result, "<script>", "")
	result = strings.ReplaceAll(result, "<script compiled>", "")
	result = strings.ReplaceAll(result, "</script>", "")

	workingDir, _ := filepath.Abs(*cli.GetArgs().InputPath)

	esbuildOutput := api.Build(api.BuildOptions{
		Stdin: &api.StdinOptions{
			Contents:   result,
			Loader:     api.LoaderTS,
			ResolveDir: ".",
		},
		Format:        api.FormatESModule,
		AbsWorkingDir: workingDir,
		Sourcemap:     api.SourceMapNone,
		Bundle:        true,
	})

	bundledContent := ""
	for _, outputFile := range esbuildOutput.OutputFiles {
		bundledContent += string(outputFile.Contents)
	}

	for _, buildError := range esbuildOutput.Errors {
		errorModel := models.Error{
			Message:  buildError.Text,
			FilePath: model.FileName(),
		}

		documentErrors.AddErrors(errorModel)
	}

	return bundledContent
}

func (model *JSFile) AddContent(partialContent string) {
	model.Content += partialContent
}

// Whether to eagerly execute the JavaScript code. Eagerly executed code is
// typically reserved for JavaScript code that effects the layout of the page.
//
// E.g. Some JavaScript that needs to run to set a property on an element that
// is not exposed through an attribute.
func (model *JSFile) IsLazy() bool {
	return true
}

func (model *JSFile) IsCompilerOnly() bool {
	// If the script is compiled, you should use the "<script compiled>" selector
	return strings.Contains(model.Content, "compiled>")
}

// A js file index that can be used for development builds.
// This should not be used during production builds as it may lead to stale data
// being served from a cdn or browser cache.
var jsFileIndex int = 0

func (model *JSFile) OutputPath() string {
	outPath := *cli.GetArgs().OutputPath
	return fmt.Sprintf("%s/%s", outPath, model.FileName())
}

// The file hash can be used to uniquely import this file.
// Warning: If the content of this file is changed, the file name will also
// change.
func (model *JSFile) FileName() string {
	if model.memoisedFileName != "" {
		return model.memoisedFileName
	}

	// If we are in development mode, we want to optimize for build times.
	// We therefore do not compute the md5 hash of the file (which is
	// computationally expensive), and instead just use an incrementing number.
	// This is less efficient for the CDN's and browser cache, but provides a
	// quicker development environment.
	//
	// TODO: remove this && false and fix browser caching
	if !*cli.GetArgs().IsProd && false {
		jsFileIndex++
		result := fmt.Sprintf("%d.js", jsFileIndex)

		model.memoisedFileName = result
		return result
	}

	hash := md5.Sum([]byte(model.Content))

	// Only return the first 8 characters to prevent overly long file names.
	// Warning: This greatly increases the probability of a hash collision.
	// TODO: Keep track of all the used hashes
	fileHash := hex.EncodeToString(hash[:8])
	result := fmt.Sprintf("%s.js", fileHash)

	model.memoisedFileName = result
	return result
}

func (model *JSFile) Format() {}
