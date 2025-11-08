package nodes

import (
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/html"
	"hudson-newey/2web/src/content/javascript"
)

func NewMarkupTextNode(lexNode *lexer.V2LexNode) *markupTextNode {
	return &markupTextNode{
		content: lexNode.Content,
	}
}

type markupTextNode struct {
	content string
}

func (model *markupTextNode) Type() string {
	return "MarkupTextNode"
}

func (model *markupTextNode) Tokens() []lexerTokens.LexToken {
	return []lexerTokens.LexToken{lexerTokens.TextContent}
}

func (model *markupTextNode) HtmlContent() *html.HTMLFile {
	return html.NewHtmlFile(model.content)
}

func (model *markupTextNode) JsContent() *javascript.JSFile {
	return javascript.NewJsFile()
}

func (model *markupTextNode) CssContent() *css.CSSFile {
	return css.NewCssFile()
}
