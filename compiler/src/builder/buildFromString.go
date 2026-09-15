package builder

import (
	"fmt"

	"hudson-newey/2web/src/cli"
	preprocessor "hudson-newey/2web/src/compiler/1-preprocessor"
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	validator "hudson-newey/2web/src/compiler/3-validator"
	parser "hudson-newey/2web/src/compiler/4-parser"
	"hudson-newey/2web/src/compiler/4-parser/grammar"
	"hudson-newey/2web/src/compiler/4-parser/nodes"
	templating "hudson-newey/2web/src/compiler/5-templating"
	"hudson-newey/2web/src/compiler/io/reader"
	"hudson-newey/2web/src/content/document/devtools"
	"hudson-newey/2web/src/content/document/documentErrors"
	"hudson-newey/2web/src/content/page"
	"hudson-newey/2web/src/models"
)

// buildFromString compiles a page from its preprocessed source.
//
// Every error produced while compiling the page (lexing, parsing, validation,
// and reactivity compilation) is funneled through the document errors pipeline
// and rendered into the page's error overlay, so that a broken page never
// silently compiles and never kills the whole build.
func buildFromString(inputPath string, data string, isFullPage bool) (compiledPage page.Page, success bool) {
	args := cli.GetArgs()

	parseContext := nodes.NewParseContext(inputPath)

	// A panic anywhere in the pipeline (e.g. from a node constructor edge case
	// that we haven't guarded) must not kill the whole build. Recovering here
	// turns the panic into a compiler error overlay for this page so that the
	// user sees where things went wrong and the remaining pages still compile.
	defer func() {
		if recovered := recover(); recovered != nil {
			internalError := models.NewError(
				fmt.Sprintf("internal compiler error: %v", recovered),
				inputPath,
				lexer.StartingPosition,
			)

			parseContext.Adopt([]*models.Error{&internalError})
			documentErrors.AddErrors(&internalError)

			compiledPage = page.NewPage()
			compiledPage.Html.Content = data
			compiledPage.Html.Content = documentErrors.InjectErrors(
				compiledPage.Html.Content,
				parseContext.Errors(),
			)

			success = false
		}
	}()

	// 1. Preprocess
	preprocessorResult := preprocessor.ProcessStaticSite(inputPath, data, isFullPage)

	// 2. Lex
	contentReader := reader.NewReader(inputPath, preprocessorResult)
	lexInstance := lexer.NewLexer(contentReader)
	lexStructure := lexInstance.Execute()
	if args.VerboseLexer {
		lexer.PrintVerboseLexer(lexStructure)
	}

	// Lexer level errors (e.g. source read failures) are attributed to the
	// page through the parse context so that they render in the page overlay.
	parseContext.Adopt(lexInstance.Errors)

	// 3. Validate
	isValid, compilerErrors := validator.IsValid(lexStructure)
	if !isValid {
		documentErrors.AddErrors(compilerErrors...)
		parseContext.Adopt(compilerErrors)
	}

	// 4. Create AST (parser)
	ast := parser.CreateAst(lexStructure, grammar.TextRules, true, parseContext)
	if args.VerboseAst {
		parser.PrintVerboseParser(ast)
	}

	// 5. Template (write result)
	compiledPage = templating.Compile(inputPath, ast)

	// Every error collected while compiling this page (parse errors, reactivity
	// errors, etc.) is rendered into the page's error overlay.
	pageErrors := compiledPage.Errors.Errors
	allPageErrors := append(parseContext.Errors(), pageErrors...)

	isErrorFree := len(allPageErrors) == 0
	if !isErrorFree {
		compiledPage.Html.Content = documentErrors.InjectErrors(
			compiledPage.Html.Content,
			allPageErrors,
		)
	}

	if args.HasDevTools {
		compiledPage.Html.Content = devtools.InjectDevTools(compiledPage.Html.Content)
	}

	return compiledPage, isErrorFree
}
