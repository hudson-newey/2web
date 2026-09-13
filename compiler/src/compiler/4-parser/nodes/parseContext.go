package nodes

import (
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/content/document/documentErrors"
	"hudson-newey/2web/src/models"
	"strings"
)

// ParseContext carries the state of the file being parsed into the grammar
// constructors so that syntax errors can be attributed to the file and the
// position in the file that caused them.
//
// Before this, malformed user syntax panicked inside the node constructors,
// which killed the entire build instead of reporting the error.
type ParseContext struct {
	FilePath string

	errors []*models.Error
}

func NewParseContext(filePath string) *ParseContext {
	return &ParseContext{FilePath: filePath}
}

// ReportError records a parse error both on the context (so that the page
// being compiled can render it in the error overlay) and in the global
// document error list (so that the compiler exit code and terminal output
// reflect it).
func (c *ParseContext) ReportError(message string, position lexer.Position) {
	errorModel := models.NewError(message, c.FilePath, position)

	c.errors = append(c.errors, &errorModel)
	documentErrors.AddErrors(&errorModel)
}

// Adopt merges externally produced errors (e.g. from the lexer) into the
// context.
func (c *ParseContext) Adopt(errors []*models.Error) {
	c.errors = append(c.errors, errors...)
}

// Errors returns the errors collected for this file.
func (c *ParseContext) Errors() []*models.Error {
	return c.errors
}

// DegradedNode reports a syntax error and returns a node that renders the raw
// source text that the grammar matched.
//
// Rendering the raw source keeps the user's (broken) content visible in the
// page while the error overlay explains what is wrong with it.
func (c *ParseContext) DegradedNode(message string, matched []*lexer.V2LexNode) *RawTextNode {
	c.ReportError(message, positionOf(matched))

	return NewRawTextNode(rawSourceText(matched))
}

// positionOf returns the position of the first lexer node in the slice.
func positionOf(lexNodes []*lexer.V2LexNode) lexer.Position {
	if len(lexNodes) == 0 {
		return lexer.StartingPosition
	}

	return lexNodes[0].Pos
}

// rawSourceText reconstructs the raw source of the consumed tokens.
func rawSourceText(lexNodes []*lexer.V2LexNode) string {
	var raw strings.Builder
	for _, lexNode := range lexNodes {
		raw.WriteString(lexNode.Content)
	}

	return raw.String()
}
