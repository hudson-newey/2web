package ast

import (
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/html"
	"hudson-newey/2web/src/content/javascript"
	"hudson-newey/2web/src/content/page"
	twoscript "hudson-newey/2web/src/content/twoScript"
)

type nodeType = string

type Node interface {
	Type() nodeType

	Content(page *page.Page) NodeContent

	Children() AbstractSyntaxTree
	AddChild(child Node)
	RemoveChild(child Node)
}

type NodeContent struct {
	// Overwrites all html of the existing page
	HtmlContent      *html.HTMLFile
	JsContent        *javascript.JSFile
	CssContent       *css.CSSFile
	TwoScriptContent *twoscript.TwoScriptFile
}
