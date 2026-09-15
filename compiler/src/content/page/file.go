package page

import (
	"fmt"
	"hudson-newey/2web/src/content"
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/document"
	"hudson-newey/2web/src/content/html"
	"hudson-newey/2web/src/content/javascript"
	twoscript "hudson-newey/2web/src/content/twoScript"
	"strings"
)

func NewPage() Page {
	return Page{
		Html:       html.FromContent(""),
		TwoScript:  []*twoscript.TwoScriptFile{},
		JavaScript: []*javascript.JSFile{},
		Css:        []*css.CSSFile{},
		Assets:     []*content.BinaryFile{},
		Errors:     &PageErrors{},
	}
}

type Page struct {
	InputPath  string
	Html       *html.HTMLFile
	TwoScript  []*twoscript.TwoScriptFile
	JavaScript []*javascript.JSFile
	Css        []*css.CSSFile
	Assets     []*content.BinaryFile
	Errors     *PageErrors

	// Ids allocates the runtime identifiers (DOM selectors, variable and
	// function names) that wire this page's compiled output together. Each page
	// owns an allocator so that the compiled output doesn't depend on which
	// other pages were compiled before it.
	Ids javascript.IdAllocator

	// reactiveRuntime accumulates the compiled reactive variable wiring of the
	// page. Every reactive variable contributes a chunk, and the chunks are
	// emitted into a single shared script (in dependency order) so that
	// computed variables can reference the runtime variables they are computed
	// from.
	ReactiveRuntime []ReactiveRuntimeChunk
}

// ReactiveRuntimeChunk is the compiled wiring of a single reactive variable.
//
// Variable is the selector of the variable the chunk belongs to and
// Dependencies lists the selectors of the variables whose chunks must be
// declared before this one (the variables this one is computed from).
type ReactiveRuntimeChunk struct {
	Variable     string
	Dependencies []string
	Content      string
}

// AppendReactiveRuntime adds a compiled reactive variable chunk to the page's
// shared reactive runtime.
func (model *Page) AppendReactiveRuntime(chunk ReactiveRuntimeChunk) {
	model.ReactiveRuntime = append(model.ReactiveRuntime, chunk)
}

// EmitReactiveRuntime flushes the accumulated reactive runtime chunks into a
// single script.
//
// The chunks are emitted in dependency order (a chunk that another chunk is
// computed from is declared first) so that computed variables can reference
// the runtime variables they are computed from. Cyclic dependencies are
// emitted in compilation order after a warning; they are reported as compiler
// errors elsewhere.
func (model *Page) EmitReactiveRuntime() {
	if len(model.ReactiveRuntime) == 0 {
		return
	}

	ordered := topologicalReactiveChunks(model.ReactiveRuntime)

	content := ""
	for _, chunk := range ordered {
		content += chunk.Content + "\n"
	}

	model.AddScript(javascript.FromGeneratedContent(content))
	model.ReactiveRuntime = nil
}

// topologicalReactiveChunks orders reactive runtime chunks so that every chunk
// appears after the chunks it depends on.
func topologicalReactiveChunks(chunks []ReactiveRuntimeChunk) []ReactiveRuntimeChunk {
	ordered := make([]ReactiveRuntimeChunk, 0, len(chunks))
	emitted := map[string]bool{}

	// A chunk is ready when none of its dependencies are pending.
	for len(ordered) < len(chunks) {
		progressed := false

		for _, chunk := range chunks {
			if emitted[chunk.Variable] {
				continue
			}

			ready := true
			for _, dependency := range chunk.Dependencies {
				if dependency == chunk.Variable {
					continue
				}

				if !emitted[dependency] {
					// Only wait for dependencies that actually have a chunk.
					// Dependencies without one (e.g. static variables) are
					// resolved at compile time.
					for _, candidate := range chunks {
						if candidate.Variable == dependency {
							ready = false
							break
						}
					}
				}
			}

			if ready {
				ordered = append(ordered, chunk)
				emitted[chunk.Variable] = true
				progressed = true
			}
		}

		if !progressed {
			// The remaining chunks form a dependency cycle. Emit them in
			// compilation order; the cycle is reported as a compiler error
			// elsewhere.
			for _, chunk := range chunks {
				if !emitted[chunk.Variable] {
					ordered = append(ordered, chunk)
					emitted[chunk.Variable] = true
				}
			}

			break
		}
	}

	return ordered
}

func (model *Page) SetContent(content string) {
	model.Html.Content = content
}

func (model *Page) SetHtmlContent(htmlFile *html.HTMLFile) {
	model.Html = htmlFile
}

func (model *Page) AddTwoScript(tsFile *twoscript.TwoScriptFile) {
	model.TwoScript = append(model.TwoScript, tsFile)
}

func (model *Page) AddScript(jsFile *javascript.JSFile) {
	if strings.TrimSpace(jsFile.Content) == "" {
		return
	}

	model.JavaScript = append(model.JavaScript, jsFile)

	// Adds a "<script src=></script>" tag to the html document to load the js file
	// If the JavaScript is not lazy loaded, we want to eagerly evaluate it by not
	// using the "async" keyword.
	// Note that this will delay the initial page load times because all of the
	// JavaScript will be blocking.
	//
	// Note that we include new line character after every import to make them
	// easier to read in a development environment.
	// These new line characters will be removed when building for production.
	injectedContent := ""
	if jsFile.IsLazy() {
		// type="module" is automatically deferred.
		injectedContent = fmt.Sprintf(`<script type="module" src="%s"></script>%s`, jsFile.FileName(), "\n")
	} else {
		// Note that the script is still deferred because it is using type="module"
		// I have purposely done this so that the programmer doesn't have to deal
		// with differing import formats when using eager/lazy scripting.
		// Eager/lazy loading should be an implementation detail transparent to the
		// programmer.
		// This also allows the user to reference DOM elements in their scripts
		// without any iife or any other related hackery.
		injectedContent = fmt.Sprintf(`<script type="module" src="%s"></script>%s`, jsFile.FileName(), "\n")
	}

	// While normal (non-module type) scripts block DOM AST construction, because
	// 2Web uses async + module type scripts, execution is deferred until the DOM
	// has been constructed.
	// Meaning that it is actually more beneficial to inject scripts into the top
	// of the head, so that they can be discovered and start fetching sooner.
	model.Html.Content = document.InjectContent(model.Html.Content, injectedContent, document.HeadTop)
}

func (model *Page) AddStyle(cssFile *css.CSSFile) {
	if strings.TrimSpace(cssFile.Content) == "" {
		return
	}

	model.Css = append(model.Css, cssFile)

	// Adds a "<link>" tag to the html document to load the css file
	injectedContent := fmt.Sprintf("<link rel=\"stylesheet\" href=\"%s\" />\n", cssFile.FileName())

	// Always inject css styles into the top of the head element so that they can
	// be discovered as soon as possible, to begin parsing
	model.Html.Content = document.InjectContent(model.Html.Content, injectedContent, document.HeadTop)
}

// Adds an asset (binary file) to the page model.
// This is an escape hatch for adding binary or unrecognized formats to the
// output assets, and a smell that the compiler cannot correctly handle the file
// type and perform necessary optimizations/transpilation/etc.
func (model *Page) AddAsset(binaryFile *content.BinaryFile) {
	model.Assets = append(model.Assets, binaryFile)
}

func (model *Page) Format() {
	if model.Html != nil {
		model.Html.Format()
	}

	for _, jsFile := range model.JavaScript {
		jsFile.Format()
	}

	for _, cssFile := range model.Css {
		cssFile.Format()
	}
}
