package lexer

import (
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
	"strings"
)

// When inside the first starting angle bracket (<) and up until (and including)
// the closing angle bracket (>).
func elementLexer(model *Lexer) (V2LexNode, LexFunc) {
	cases := lexDefMap{
		"!doctype": {token: lexeme.Doctype, next: elementLexer},
		"script compiled": {token: lexeme.CompiledScriptStartTag, next: inlineCompiledScriptTagLexer},
		"script":          {token: lexeme.ScriptStartTag, next: inlineScriptTagLexer},
		"style":           {token: lexeme.StyleStartTag, next: inlineStyleTagLexer},
		"code":            {token: lexeme.CodeStartTag, next: inlineCodeTagLexer},
	}

	// "style", "script", and "code" are also common attribute names
	// (e.g. style="color: red"). HTML tokenizers treat <style="x"> as an
	// unknown element rather than a style tag, so these matchers must not fire
	// when the tag name is immediately followed by an assignment. Without this,
	// a style attribute would switch the lexer into the style content state and
	// everything up to the next </style> in the document would be swallowed as
	// style content.
	//
	// Note that attributes without a value (e.g. <span data-style>) are still
	// treated as tag openings; disambiguating those requires full parser level
	// context that the lexer doesn't have.
	for _, tagName := range []string{"script", "style", "code"} {
		if attributeAssignment(model, tagName) {
			delete(cases, tagName)
		}
	}

	cases = withAttributes(cases, elementLexer)
	cases = withStrings(cases, elementLexer)

	return lexerFactory(cases, element)(model)
}

// attributeAssignment returns whether the input at the current position reads
// like an attribute assignment rather than a tag name.
// e.g. style="color: red" or data-code="123"
func attributeAssignment(model *Lexer, tagName string) bool {
	peeked := strings.ToLower(model.peek(len(tagName) + 1))

	if len(peeked) <= len(tagName) {
		return false
	}

	if !strings.HasPrefix(peeked, strings.ToLower(tagName)) {
		return false
	}

	return peeked[len(tagName)] == '='
}
