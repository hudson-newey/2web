package nodes

import (
	"errors"
	"fmt"
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
	"hudson-newey/2web/src/compiler/4-parser/scanners"
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/html"
	"hudson-newey/2web/src/content/javascript"
	"hudson-newey/2web/src/content/page"
	twoscript "hudson-newey/2web/src/content/twoScript"
	"hudson-newey/2web/src/models"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// import importName from "importPath";
func NewscriptImportNode(lexNodes []*lexer.V2LexNode) *scriptImportNode {
	importNameNode, err := scanners.NthToken(lexNodes, lexeme.CompiledScriptSource, 1)
	importNameNode.Content = strings.TrimSpace(importNameNode.Content)
	if err != nil {
		panic(err)
	}

	importPathNode, err := scanners.NthToken(lexNodes, lexeme.CompiledScriptSource, 2)
	quoteRe := regexp.MustCompile(`"(.*?)"`)
	importPathNode.Content = quoteRe.FindAllStringSubmatch(importPathNode.Content, -1)[0][1]
	if err != nil {
		panic(err)
	}

	return &scriptImportNode{
		importName: importNameNode.Content,
		importPath: importPathNode.Content,
	}
}

type scriptImportNode struct {
	// What was the alias that the import was given?
	importName string

	// What is the path of the import?
	importPath string

	children AbstractSyntaxTree
}

func (m *scriptImportNode) Type() string {
	return "importNode"
}

func (m *scriptImportNode) Children() AbstractSyntaxTree {
	return m.children
}

func (m *scriptImportNode) Content(page *page.Page, _ast AbstractSyntaxTree) NodeContent {
	// Check that the imported path really exists.
	hostDirectory := filepath.Dir(page.InputPath)
	componentPath := filepath.Join(hostDirectory, m.importPath)
	componentContent, err := os.ReadFile(componentPath)
	if errors.Is(err, os.ErrNotExist) {
		msg := fmt.Sprintf("error importing file. Could not find file: '%s'", m.importPath)
		pageErr := models.NewError(msg, page.InputPath, lexer.Position{})
		page.Errors.AddError(&pageErr)
	}

	// Replace all of the selectors with the components content.
	selector := fmt.Sprintf("<%s />", m.importName)
	HtmlContent := strings.ReplaceAll(page.Html.Content, selector, string(componentContent))

	return NodeContent{
		HtmlContent:      html.FromContent(HtmlContent),
		JsContent:        javascript.NewJsFile(),
		CssContent:       css.NewCssFile(),
		TwoScriptContent: twoscript.NewTwoScriptFile(),
	}
}

func (m *scriptImportNode) AddChild(child Node) {
	m.children = append(m.children, child)
}

func (m *scriptImportNode) RemoveChild(child Node) {
	for i, c := range m.children {
		if c == child {
			m.children = append(m.children[:i], m.children[i+1:]...)
			return
		}
	}
}
