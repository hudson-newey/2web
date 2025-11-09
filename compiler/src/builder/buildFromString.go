package builder

import (
	"hudson-newey/2web/src/cli"
	preprocessor "hudson-newey/2web/src/compiler/1-preprocessor"
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	validator "hudson-newey/2web/src/compiler/3-validator"
	parser "hudson-newey/2web/src/compiler/4-parser"
	templating "hudson-newey/2web/src/compiler/5-templating"
	"hudson-newey/2web/src/compiler/io/reader"
	"hudson-newey/2web/src/content/document/devtools"
	"hudson-newey/2web/src/content/document/documentErrors"
	"hudson-newey/2web/src/content/page"
)

func buildFromString(inputPath string, data string, isFullPage bool) (page.Page, bool) {
	args := cli.GetArgs()

	// 1. Preprocess
	preprocessorResult := preprocessor.ProcessStaticSite(inputPath, data, isFullPage)

	// 2. Lex
	contentReader := reader.NewReader(inputPath, preprocessorResult)
	lexInstance := lexer.NewLexer(contentReader)
	lexStructure := lexInstance.Execute()
	if *args.VerboseLexer {
		lexer.PrintVerboseLexer(lexStructure)
	}

	// 3. Validate
	isValid, compilerErrors := validator.IsValid(lexStructure)
	if !isValid {
		documentErrors.AddErrors(compilerErrors...)
	}

	// 4. Create AST (parser)
	ast := parser.CreateAst(lexStructure)
	if *args.VerboseAst {
		parser.PrintVerboseParser(ast)
	}

	// 5. Template (write result)
	compiledPage := templating.Compile(inputPath, ast)

	isErrorFree := documentErrors.IsPageErrorFree()
	if !isErrorFree {
		compiledPage.Html.Content = documentErrors.InjectErrors(compiledPage.Html.Content)
		documentErrors.ResetPageErrors()
	}

	if *args.HasDevTools {
		compiledPage.Html.Content = devtools.InjectDevTools(compiledPage.Html.Content)
	}

	return compiledPage, isErrorFree
}
