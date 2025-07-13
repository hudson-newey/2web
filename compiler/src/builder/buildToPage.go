package builder

import (
	"hudson-newey/2web/src/cli"
	preprocessor "hudson-newey/2web/src/compiler/1-preprocessor"
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	validator "hudson-newey/2web/src/compiler/3-validator"
	parser "hudson-newey/2web/src/compiler/4-parser"
	templating "hudson-newey/2web/src/compiler/5-templating"
	"hudson-newey/2web/src/content/document/devtools"
	"hudson-newey/2web/src/content/document/documentErrors"
	"hudson-newey/2web/src/content/page"
	"hudson-newey/2web/src/optimizer"
)

func buildToPage(inputPath string) (page.Page, bool) {
	args := cli.GetArgs()

	data := getContent(inputPath)

	// 1. Preprocess
	ssgResult := preprocessor.ProcessStaticSite(inputPath, data)
	pageModel := templating.BuildPage(ssgResult)

	// 2. Lex
	lexInstance := lexer.NewLexer(pageModel.Html.Reader())
	lexStructure := lexInstance.Execute()

	// 3. Validate
	isValid, compilerErrors := validator.IsValid(lexStructure)
	if !isValid {
		documentErrors.AddErrors(compilerErrors...)
	}

	// 4. Create AST (parser)
	ast := parser.CreateAst(lexStructure)

	// 5. Template (write result)
	compiledPage := templating.Compile(inputPath, pageModel, ast)

	isErrorFree := documentErrors.IsPageErrorFree()
	if !isErrorFree {
		compiledPage.Html.Content = documentErrors.InjectErrors(compiledPage.Html.Content)
		documentErrors.ResetPageErrors()
	}

	if *args.HasDevTools {
		compiledPage.Html.Content = devtools.InjectDevTools(compiledPage.Html.Content)
	}

	// We always optimize last so that even the injected content is optimized.
	if *args.IsProd {
		compiledPage = optimizer.OptimizePage(compiledPage)
	}

	return compiledPage, isErrorFree
}
