package nodes

import (
	"strings"

	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
)

// Functions declared inside a <script compiled> block are passed through to
// the page's compiled JavaScript as normal functions, with two reactive
// extensions:
//
//  1. Assignments to reactive variables inside the function body update the
//     page (they are rewritten into the same "assign the runtime variable and
//     run its update cascade" shape that event listeners compile into).
//  2. Server functions that were imported from a .server.ts script are in
//     scope, so they can be called like normal async functions.
//
// The block source is reconstructed from the lexed tokens (the lexer has no
// brace context, so a '$name = value;' statement inside a function body is
// lexed exactly like a top level declaration). Functions are extracted with a
// string and comment aware brace matcher so that the rest of the block (the
// reactive variable declarations) is untouched.

type sourcePosition struct {
	row int
	col int
}

type reactiveReferenceKind int

const (
	// A read of a reactive variable, e.g. "return $count + 1".
	refRead reactiveReferenceKind = iota

	// An assignment to a reactive variable, e.g. "$count = $count + 1;".
	refAssign

	// A compound assignment, e.g. "$count += 2;".
	refCompoundAssign

	// An increment/decrement statement, e.g. "$count++;".
	refIncrement
)

// reactiveReference is a single reference to a reactive variable inside a
// compiled function body.
type reactiveReference struct {
	selector string
	kind     reactiveReferenceKind

	// The offset of the "$name" selector within the function source.
	selectorStart int

	// The span of the source that the reference rewrites. For reads this is
	// just the "$name" selector; for assignments it extends through the
	// statement's terminating semicolon (when there is one).
	spanStart int
	spanEnd   int

	// The span of the assignment operator ("=", "+=", "++", ...) within the
	// rewritten span. Zero valued for reads.
	opStart int
	opEnd   int

	// Whether the rewritten span includes its terminating semicolon.
	hasSemicolon bool

	// The absolute position of the selector in the source file.
	position sourcePosition
}

// compiledFunction is a function declaration extracted from a compiled script
// block.
type compiledFunction struct {
	name       string
	source     string
	start      sourcePosition
	end        sourcePosition
	references []reactiveReference
}

// sourceToken is a content token of a compiled script block, mapped into the
// reconstructed block source.
type sourceToken struct {
	content string
	start   sourcePosition

	// The span of the token's content within the reconstructed source.
	outStart int
	outEnd   int
}

// scriptBlockTokens filters the grammar's structural tokens (the script tags)
// out of the matched token subset, leaving the block's content tokens.
func scriptBlockTokens(lexNodes []*lexer.V2LexNode) []*lexer.V2LexNode {
	content := []*lexer.V2LexNode{}

	for _, lexNode := range lexNodes {
		switch lexNode.Token {
		case lexeme.LessAngle, lexeme.CompiledScriptStartTag, lexeme.GreaterAngle, lexeme.ScriptEndTag:
			continue
		}

		content = append(content, lexNode)
	}

	return content
}

// reconstructBlockSource rebuilds the block's source text from its content
// tokens.
//
// The lexer drops the newlines between tokens (whitespace immediately after a
// statement's semicolon isn't part of any token), so the token positions are
// used to re-insert them: whenever the next token starts on a later row than
// the reconstructed source has reached, the missing newlines are written back
// in. This keeps the emitted functions faithful to the user's source.
//
// It also returns the token mapping into the reconstructed source so that
// match offsets can be translated back into absolute file positions.
func reconstructBlockSource(lexNodes []*lexer.V2LexNode) (string, []sourceToken) {
	var source strings.Builder
	tokens := []sourceToken{}

	for _, lexNode := range scriptBlockTokens(lexNodes) {
		start := sourcePosition{row: lexNode.Pos.Row, col: lexNode.Pos.Col}

		// Re-insert the newlines the lexer skipped between this token and the
		// previous one.
		if len(tokens) > 0 {
			previous := tokens[len(tokens)-1]
			expected := advancePosition(previous.start, previous.content)

			for missing := start.row - expected.row; missing > 0; missing-- {
				source.WriteByte('\n')
			}
		}

		outStart := source.Len()
		source.WriteString(lexNode.Content)

		tokens = append(tokens, sourceToken{
			content:  lexNode.Content,
			start:    start,
			outStart: outStart,
			outEnd:   source.Len(),
		})
	}

	return source.String(), tokens
}

// advancePosition returns the source position after consuming the given text
// from the given position. Columns are counted in runes, matching the lexer.
func advancePosition(pos sourcePosition, text string) sourcePosition {
	newlines := strings.Count(text, "\n")

	if newlines > 0 {
		pos.row += newlines
		pos.col = runeLen(text[strings.LastIndexByte(text, '\n')+1:]) + 1
		return pos
	}

	pos.col += runeLen(text)
	return pos
}

func runeLen(text string) int {
	count := 0
	for range text {
		count++
	}

	return count
}

// positionAt translates an offset in the reconstructed source into an
// absolute file position.
func positionAt(offset int, source string, tokens []sourceToken) sourcePosition {
	for _, token := range tokens {
		if offset < token.outStart || offset >= token.outEnd {
			continue
		}

		return advancePosition(token.start, token.content[:offset-token.outStart])
	}

	// The offset falls between tokens (in re-inserted whitespace). Walk
	// backwards to the nearest token that ends before the offset and advance
	// over the gap.
	for i := len(tokens) - 1; i >= 0; i-- {
		if tokens[i].outEnd <= offset {
			return advancePosition(tokens[i].start, tokens[i].content)
		}
	}

	return sourcePosition{row: 1, col: 1}
}

// extractCompiledFunctions extracts the top level function declarations from
// a compiled script block's lexed tokens.
func extractCompiledFunctions(lexNodes []*lexer.V2LexNode) []compiledFunction {
	source, tokens := reconstructBlockSource(lexNodes)
	locate := func(offset int) sourcePosition {
		return positionAt(offset, source, tokens)
	}

	functions := []compiledFunction{}
	depth := 0

	for i := 0; i < len(source); {
		char := source[i]

		switch {
		case char == '\'' || char == '"':
			i = skipQuotedString(source, i, char)

		case char == '`':
			i = skipTemplateLiteral(source, i)

		case strings.HasPrefix(source[i:], "//"):
			i = skipUntil(source, i, '\n')

		case strings.HasPrefix(source[i:], "/*"):
			i = skipUntil(source, i+2, '*')
			if i < len(source) {
				i++
			}

		case char == '{':
			depth++
			i++

		case char == '}':
			depth--
			i++

		case isIdentifierByte(char):
			word, wordEnd := readWord(source, i)

			// Only top level function declarations are extracted. Nested
			// declarations stay verbatim inside their enclosing function.
			if depth != 0 || (word != "function" && word != "async") {
				i = wordEnd
				continue
			}

			declarationStart := i
			keywordEnd := wordEnd

			if word == "async" {
				keywordEnd = skipSpaceAndNewlines(source, keywordEnd)

				nextWord, nextWordEnd := readWord(source, keywordEnd)
				if nextWord != "function" {
					i = wordEnd
					continue
				}

				keywordEnd = nextWordEnd
			}

			function, ok := matchFunctionDeclaration(source, declarationStart, keywordEnd)
			if !ok {
				i = keywordEnd
				continue
			}

			function.start = locate(declarationStart)
			function.end = locate(declarationStart + function.spanEnd() - 1)
			function.references = scanReactiveReferences(function.source, func(offset int) sourcePosition {
				return locate(declarationStart + offset)
			})

			functions = append(functions, function)
			i = declarationStart + len(function.source)

		default:
			i++
		}
	}

	return functions
}

// matchFunctionDeclaration matches "function name(...)" (the keyword span is
// already consumed) and returns the full declaration source, from
// declarationStart through the closing brace.
func matchFunctionDeclaration(source string, declarationStart int, keywordEnd int) (compiledFunction, bool) {
	cursor := skipSpaceAndNewlines(source, keywordEnd)

	name, nameEnd := readWord(source, cursor)
	if name == "" {
		return compiledFunction{}, false
	}

	cursor = skipSpaceAndNewlines(source, nameEnd)
	if cursor >= len(source) || source[cursor] != '(' {
		return compiledFunction{}, false
	}

	// The parameter list.
	depth := 0
	for ; cursor < len(source); cursor++ {
		char := source[cursor]

		if char == '\'' || char == '"' {
			cursor = skipQuotedString(source, cursor, char) - 1
			continue
		}

		if char == '(' {
			depth++
		}

		if char == ')' {
			depth--
			if depth == 0 {
				break
			}
		}
	}

	// The body.
	for ; cursor < len(source); cursor++ {
		if source[cursor] == '{' {
			break
		}
	}

	if cursor >= len(source) {
		return compiledFunction{}, false
	}

	bodyEnd, ok := matchBrace(source, cursor)
	if !ok {
		return compiledFunction{}, false
	}

	return compiledFunction{name: name, source: source[declarationStart:bodyEnd]}, true
}

// spanEnd returns the length of the function's source (the exclusive end
// offset within the block source).
func (fn compiledFunction) spanEnd() int {
	return len(fn.source)
}

// scanReactiveReferences finds every reference to a reactive variable in the
// given function source: reads, assignments, compound assignments, and the
// increment/decrement shorthand.
//
// Only references in statement position are treated as assignments (a
// statement position is the start of the source, or directly after ';', '{',
// '}', ')' or an 'else'/'do' keyword); everywhere else the reactive variable
// is a plain read and is only substituted with its runtime value.
func scanReactiveReferences(source string, locate func(int) sourcePosition) []reactiveReference {
	references := []reactiveReference{}
	statementStart := true

	for i := 0; i < len(source); {
		char := source[i]

		switch {
		case char == '\'' || char == '"':
			i = skipQuotedString(source, i, char)

		case char == '`':
			i = skipTemplateLiteral(source, i)

		case strings.HasPrefix(source[i:], "//"):
			i = skipUntil(source, i, '\n')

		case strings.HasPrefix(source[i:], "/*"):
			i = skipUntil(source, i+2, '*')
			if i < len(source) {
				i++
			}

		case char == ';', char == '{', char == '}':
			statementStart = true
			i++

		case char == ')':
			statementStart = true
			i++

		case char == '$' && i+1 < len(source) && isIdentifierByte(source[i+1]):
			reference, ok := matchReactiveReference(source, i, statementStart)
			if !ok {
				// Not a reactive variable after all (e.g. "$$" or an
				// identifier continuation). Copy the character.
				statementStart = false
				i++
				continue
			}

			reference.position = locate(reference.spanStart)
			references = append(references, reference)

			if reference.kind == refRead {
				statementStart = false
				i = reference.selectorEnd()
				continue
			}

			// Assignments consume their whole statement.
			statementStart = true
			i = reference.spanEnd

		case isIdentifierByte(char):
			word, wordEnd := readWord(source, i)
			if word != "else" && word != "do" {
				statementStart = false
			}

			i = wordEnd

		default:
			if char > ' ' {
				statementStart = false
			}

			i++
		}
	}

	return references
}

// matchReactiveReference matches "$name" at the given offset and classifies
// the reference by the operator that follows it.
func matchReactiveReference(source string, offset int, statementStart bool) (reactiveReference, bool) {
	_, selectorEnd := readWord(source, offset+1)
	reference := reactiveReference{
		spanStart:     offset,
		selectorStart: offset,
	}

	reference.setSelector(source[offset:selectorEnd])

	cursor := skipSpaces(source, selectorEnd)

	rest := source[cursor:]
	next := byte(0)
	if cursor < len(source) {
		next = source[cursor]
	}

	switch {
	case strings.HasPrefix(rest, "++") || strings.HasPrefix(rest, "--"):
		if !statementStart {
			// An increment inside an expression (e.g. "arr[$count++]") is a
			// read: the selector is substituted and the increment stays
			// verbatim.
			return reference, true
		}

		reference.kind = refIncrement
		reference.opStart = cursor
		reference.opEnd = cursor + 2
		reference.spanEnd = reference.opEnd

		if cursor+2 < len(source) && source[cursor+2] == ';' {
			reference.hasSemicolon = true
			reference.spanEnd++
		}

		return reference, true

	case next == '=' && !strings.HasPrefix(rest, "==") && !strings.HasPrefix(rest, "=>"):
		return matchAssignmentReference(source, reference, selectorEnd, cursor, "=")

	case strings.HasPrefix(rest, "+=") || strings.HasPrefix(rest, "-=") ||
		strings.HasPrefix(rest, "*=") || strings.HasPrefix(rest, "/=") ||
		strings.HasPrefix(rest, "%="):
		if !statementStart {
			return reference, true
		}

		return matchAssignmentReference(source, reference, selectorEnd, cursor, source[cursor:cursor+2])
	}

	return reference, true
}

// matchAssignmentReference completes an assignment reference by finding the
// end of its value expression.
func matchAssignmentReference(source string, reference reactiveReference, selectorEnd int, opStart int, op string) (reactiveReference, bool) {
	reference.kind = refAssign
	if op != "=" {
		reference.kind = refCompoundAssign
	}

	reference.opStart = opStart
	reference.opEnd = opStart + len(op)

	depth := 0
	cursor := reference.opEnd

	for ; cursor < len(source); cursor++ {
		char := source[cursor]

		switch {
		case char == '\'' || char == '"':
			cursor = skipQuotedString(source, cursor, char) - 1

		case char == '`':
			cursor = skipTemplateLiteral(source, cursor) - 1

		case char == '(' || char == '[':
			depth++

		case char == ')' || char == ']':
			depth--

		case char == '{':
			// A '{' at expression depth ends the statement (e.g. the body of
			// "if (...) $count = 1 { ... }" style code is malformed anyway);
			// treat it as the end of the expression.
			if depth == 0 {
				reference.spanEnd = cursor
				return reference, true
			}

			depth++

		case char == '}':
			if depth == 0 {
				reference.spanEnd = cursor
				return reference, true
			}

			depth--

		case char == ';' && depth == 0:
			reference.hasSemicolon = true
			reference.spanEnd = cursor + 1
			return reference, true
		}
	}

	reference.spanEnd = len(source)
	return reference, true
}

// selectorEnd returns the exclusive end offset of the reference's selector.
func (ref *reactiveReference) selectorEnd() int {
	return ref.selectorStart + len(ref.selector)
}

func (ref *reactiveReference) setSelector(selector string) {
	ref.selector = selector
}

// ---------------------------------------------------------------------------
// Source scanning helpers
// ---------------------------------------------------------------------------

// readWord reads the identifier at the given offset. Returns the word (empty
// when there is no identifier) and the offset after it.
func readWord(source string, offset int) (string, int) {
	end := offset

	for end < len(source) && isIdentifierByte(source[end]) {
		end++
	}

	return source[offset:end], end
}

func skipSpaceAndNewlines(source string, offset int) int {
	for offset < len(source) && (source[offset] == ' ' || source[offset] == '\t' || source[offset] == '\n' || source[offset] == '\r') {
		offset++
	}

	return offset
}

func skipSpaces(source string, offset int) int {
	for offset < len(source) && (source[offset] == ' ' || source[offset] == '\t') {
		offset++
	}

	return offset
}

// skipQuotedString returns the offset after the closing quote of the string
// that starts at the given offset.
func skipQuotedString(source string, offset int, quote byte) int {
	offset++

	for offset < len(source) {
		if source[offset] == '\\' {
			offset += 2
			continue
		}

		if source[offset] == quote {
			return offset + 1
		}

		offset++
	}

	return offset
}

// skipTemplateLiteral returns the offset after the closing backtick of the
// template literal that starts at the given offset.
func skipTemplateLiteral(source string, offset int) int {
	offset++

	for offset < len(source) {
		if source[offset] == '\\' {
			offset += 2
			continue
		}

		if source[offset] == '`' {
			return offset + 1
		}

		offset++
	}

	return offset
}

// skipUntil returns the offset after the next occurrence of the given
// character.
func skipUntil(source string, offset int, char byte) int {
	for offset < len(source) {
		current := source[offset]
		offset++

		if current == char {
			break
		}
	}

	return offset
}

// matchBrace returns the offset after the '}' that closes the '{' at the
// given offset.
func matchBrace(source string, open int) (int, bool) {
	depth := 0

	for i := open; i < len(source); {
		char := source[i]

		switch {
		case char == '\'' || char == '"':
			i = skipQuotedString(source, i, char)

		case char == '`':
			i = skipTemplateLiteral(source, i)

		case strings.HasPrefix(source[i:], "//"):
			i = skipUntil(source, i, '\n')

		case strings.HasPrefix(source[i:], "/*"):
			i = skipUntil(source, i+2, '*')
			if i < len(source) {
				i++
			}

		case char == '{':
			depth++
			i++

		case char == '}':
			depth--
			if depth == 0 {
				return i + 1, true
			}

			i++

		default:
			i++
		}
	}

	return open, false
}

// ---------------------------------------------------------------------------
// Rewriting
// ---------------------------------------------------------------------------

// rewriteCompiledFunction rewrites the reactive references in a compiled
// function against the page's reactive index:
//
//   - assignments (including compound assignments and the increment/decrement
//     shorthand) become runtime variable assignments followed by the
//     variable's update function call, so that the page updates;
//   - reads become references to the variable's runtime value (static
//     variables are inlined with their resolved value).
func rewriteCompiledFunction(fn compiledFunction, index *ReactiveIndex) string {
	rewritten := fn.source

	// The references are spliced in reverse so that the offsets of the not
	// yet rewritten references stay valid.
	for i := len(fn.references) - 1; i >= 0; i-- {
		reference := fn.references[i]
		runtimeName := index.RuntimeVariableName(reference.selector)

		if reference.kind == refRead {
			if runtimeName == "" {
				resolved, ok := resolveStaticSelector(reference.selector, index)
				if !ok {
					continue
				}

				rewritten = rewritten[:reference.spanStart] + resolved + rewritten[reference.selectorEnd():]
				continue
			}

			rewritten = rewritten[:reference.spanStart] + runtimeName + rewritten[reference.selectorEnd():]
			continue
		}

		// An assignment to a variable without a runtime representation can't
		// trigger updates; fall back to substituting the selector with the
		// variable's resolved value so that the function stays executable.
		if runtimeName == "" {
			continue
		}

		handlerName := index.HandlerName(reference.selector)
		if handlerName == "" {
			continue
		}

		rewritten = rewritten[:reference.spanStart] + rewriteAssignment(reference, fn.source, runtimeName, handlerName, index) + rewritten[reference.spanEnd:]
	}

	return rewritten
}

// rewriteAssignment builds the replacement source for an assignment
// reference: the runtime variable assignment followed by the variable's
// update function call.
func rewriteAssignment(reference reactiveReference, source string, runtimeName string, handlerName string, index *ReactiveIndex) string {
	switch reference.kind {
	case refIncrement:
		operator := source[reference.opStart:reference.opEnd]

		return runtimeName + operator + "; " + handlerName + "(" + runtimeName + ");"

	case refCompoundAssign:
		operator := source[reference.opStart:reference.opEnd]
		expression := referenceValueExpression(reference, source, index)

		return runtimeName + " " + operator + " " + expression + "; " + handlerName + "(" + runtimeName + ");"
	}

	expression := referenceValueExpression(reference, source, index)

	return runtimeName + " = " + expression + "; " + handlerName + "(" + runtimeName + ");"
}

// referenceValueExpression returns the value expression of an assignment
// reference (everything between the operator and the statement's end) with
// the reactive variables it reads resolved to their runtime values.
func referenceValueExpression(reference reactiveReference, source string, index *ReactiveIndex) string {
	end := reference.spanEnd
	if reference.hasSemicolon {
		end--
	}

	return resolveReactiveExpression(strings.TrimSpace(source[reference.opEnd:end]), index)
}

// resolveReactiveExpression substitutes every reactive variable that an
// expression reads with its runtime value.
//
// A variable that has a runtime representation (it is assigned by events or
// functions, or computed at runtime) is replaced with its runtime variable
// name. A static variable is replaced with its resolved initial expression,
// so the expression reads the same value the rest of the page renders.
func resolveReactiveExpression(expression string, index *ReactiveIndex) string {
	for _, variable := range index.Variables {
		if !containsSelector(expression, variable.selector()) {
			continue
		}

		if runtimeName := index.RuntimeVariableName(variable.selector()); runtimeName != "" {
			expression = replaceSelector(expression, variable.selector(), runtimeName)
			continue
		}

		resolved, ok := index.ResolveExpression(variable)
		if !ok {
			// The variable is part of a cyclic dependency, which is reported
			// as a compiler error elsewhere. Leave the selector unreplaced
			// rather than emitting a broken expression.
			continue
		}

		expression = replaceSelector(expression, variable.selector(), "("+resolved+")")
	}

	return expression
}

// resolveStaticSelector returns the resolved value expression of a static
// (non runtime) reactive variable.
func resolveStaticSelector(selector string, index *ReactiveIndex) (string, bool) {
	variable := index.VariableBySelector(selector)
	if variable == nil {
		return "", false
	}

	return index.ResolveExpression(variable)
}
