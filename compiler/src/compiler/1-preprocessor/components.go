package preprocessor

import (
	"fmt"
	"html"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"hudson-newey/2web/src/content/assets"
)

// Matches static ESM style import statements.
// e.g. import Header from "components/header.component.html";
//
// Multi line imports (e.g. import {\n\tHeader\n} from "...") are supported.
var componentImportPattern = regexp.MustCompile(`(?s)import\s+([^;]+?)\s+from\s+("[^"]+"|'[^']+')\s*;?`)

// Matches a "$props().<name>" access. e.g. "$props().title"
var propsAccessPattern = regexp.MustCompile(`\$props\(\)\.([A-Za-z_$][A-Za-z0-9_$]*)`)

// Matches the "$props()" object form.
var propsObjectPattern = regexp.MustCompile(`\$props\(\)`)

// Matches an instance parameter. e.g. [title]="'My Counter'"
var instancePropPattern = regexp.MustCompile(`\[([A-Za-z_$][A-Za-z0-9_$]*)\]="([^"]*)"`)

// isInteropImport returns whether the import imports the compiler's compile
// time interop types (the virtual functions like "$props()" and "$uid()").
func isInteropImport(importPath string) bool {
	return strings.Contains(importPath, "interop.types")
}

// isConstantValue returns whether a component parameter value is a compile
// time constant (a string literal, number, or boolean), which can be inlined
// into the component's markup.
func isConstantValue(value string) bool {
	if strings.Contains(value, "$") {
		return false
	}

	if regexp.MustCompile(`^(?:'[^']*'|"[^"]*"|-?[0-9]+(?:\.[0-9]+)?|true|false|null)$`).MatchString(value) {
		return true
	}

	return false
}

// templateEscape formats a constant component parameter for inlining into a
// text output: string literals are unquoted and html escaped (matching the
// text rendering semantics of a runtime text output); numbers and booleans
// are rendered as-is.
func templateEscape(value string) string {
	unquoted := regexp.MustCompile(`^'([^']*)'$|^"([^"]*)"$`).FindStringSubmatch(value)

	if unquoted != nil {
		if unquoted[1] != "" {
			return html.EscapeString(unquoted[1])
		}

		return html.EscapeString(unquoted[2])
	}

	return html.EscapeString(value)
}

// propsObjectLiteral builds the object literal that the "$props()" object
// form evaluates to.
func propsObjectLiteral(instanceProps map[string]string) string {
	entries := []string{}
	for name, value := range instanceProps {
		entries = append(entries, name+": "+value)
	}

	if len(entries) == 0 {
		return "{}"
	}

	return "{ " + strings.Join(entries, ", ") + " }"
}

// expandComponents inlines imported components into the page source.
//
// Components are expanded in the preprocessor (instead of at template time,
// which is what the script import node does as a fallback) so that the whole
// compilation pipeline can see the component's content. This matters because
// components commonly contain <script compiled> blocks, reactive properties,
// and events of their own. Expanding at template time would inline those
// blocks as raw text into the compiled page, because the page has already been
// lexed and parsed by then.
//
// Import statements are resolved relative to the importing file, so a
// component can import another component that sits next to it.
func expandComponents(filePath string, content string, visiting map[string]bool, scopeCounter *int) string {
	importDirectory := filepath.Dir(filePath)

	for _, match := range componentImportPattern.FindAllStringSubmatch(content, -1) {
		statement := match[0]
		importName := strings.TrimSpace(match[1])
		importPath := unquoteImportPath(match[2])

		componentPath := filepath.Clean(filepath.Join(importDirectory, importPath))

		// Cycle guard: a component import chain that loops back on itself must
		// not recurse forever. The revisited component is left unexpanded.
		if visiting[componentPath] {
			continue
		}

		componentContent, err := os.ReadFile(componentPath)
		if err != nil {
			// Missing imports are left untouched so that the import node can
			// surface them as compiler errors.
			continue
		}

		// Compile time virtual functions (e.g. "$props()") are imported from
		// the compiler's interop types. They are evaluated by the compiler
		// itself, so the import statement is removed instead of being
		// resolved.
		if isInteropImport(importPath) {
			content = strings.Replace(content, statement, "", 1)
			continue
		}

		// Only markup files are components. e.g. ESM imports of .ts files are
		// bundled by esbuild and must be left alone.
		if !assets.IsMarkupFile(componentPath) {
			continue
		}

		// Components are inlined by replacing their selector
		// (e.g. <Header />). A component without a selector in the importing
		// file isn't used, so inlining it would only pollute the page with its
		// script variables.
		selectorPattern, selectorPatternErr := componentSelectorPattern(importName)
		if selectorPatternErr != nil {
			continue
		}

		if !selectorPattern.MatchString(content) {
			continue
		}

		visiting[componentPath] = true
		expandedComponent := expandComponents(componentPath, string(componentContent), visiting, scopeCounter)
		delete(visiting, componentPath)

		content = strings.Replace(content, statement, "", 1)
		content = selectorPattern.ReplaceAllStringFunc(content, func(selector string) string {
			return expandComponentInstance(selector, expandedComponent, scopeCounter)
		})
	}

	return content
}

// componentSelectorPattern matches an instance of the component with the
// given import name, including its attributes and its children.
//
// Both self closing instances and instances with children are matched:
//
//	<Counter [title]="'My Counter'" />
//	<Counter><p slot="header">Hello</p></Counter>
//
// The second capture group holds the instance's children (empty for self
// closing instances), which the component's slots render.
func componentSelectorPattern(importName string) (*regexp.Regexp, error) {
	if !regexp.MustCompile(`^[A-Za-z][A-Za-z0-9]*$`).MatchString(importName) {
		// The import name isn't a usable markup tag name (e.g. a braced esm
		// import).
		return nil, fmt.Errorf("not a component import")
	}

	name := regexp.QuoteMeta(importName)

	return regexp.MustCompile(
		`(?s)<` + name + `\b([^>]*?)(?:/>|>(.*?)</` + name + `\s*>)`,
	), nil
}

// expandComponentInstance inlines a component instance's content, substituting
// the instance's parameters for the "$props()" accesses in the component, and
// scoping the component's styles to the instance.
//
// Components without style blocks aren't scoped: their markup compiles
// unchanged.
func expandComponentInstance(selector string, componentContent string, scopeCounter *int) string {
	attributes, children := parseInstance(selector)

	instanceProps := parseInstanceProps(attributes)

	content := componentContent

	// The instance parameters replace the "$props()" accesses of the
	// component. A reactive value keeps its expression: the component
	// references the reactive variable the instance passed, so the component
	// updates when it changes. A constant value is inlined into text outputs
	// (escaped, matching the text rendering semantics of a runtime text
	// output).
	for name, value := range instanceProps {
		if !isConstantValue(value) {
			content = strings.ReplaceAll(content, "$props()."+name, value)
			continue
		}

		textOutput := regexp.MustCompile(
			`(\{\{\s*)\$props\(\)\.` + regexp.QuoteMeta(name) + `(\s*\}\})`,
		)

		content = textOutput.ReplaceAllString(content, "${1}"+templateEscape(value)+"${2}")
		content = strings.ReplaceAll(content, "$props()."+name, value)
	}

	// The remaining "$props()" accesses (properties the instance didn't pass,
	// and the object form) resolve to undefined / an object of the passed
	// values.
	content = propsObjectPattern.ReplaceAllString(content, propsObjectLiteral(instanceProps))
	content = propsAccessPattern.ReplaceAllString(content, "undefined")

	if hasStyleBlock(content) {
		*scopeCounter++
		content = scopeComponent(content, *scopeCounter)
	}

	// The instance's children fill the component's slots. This runs after the
	// props and scoping, so the slot contents are the consumer's own markup:
	// they don't receive the component's props, and they aren't tagged with
	// the component's scope attribute.
	return spliceSlots(content, children)
}

// instanceProps maps the parameter names of a component instance to the
// (raw) value expressions they were passed.
// e.g. <Counter [title]="'My Counter'" [count]="1" />.
func parseInstanceProps(attributes string) map[string]string {
	props := map[string]string{}

	for _, match := range instancePropPattern.FindAllStringSubmatch(attributes, -1) {
		props[match[1]] = strings.TrimSpace(match[2])
	}

	return props
}

func unquoteImportPath(quotedPath string) string {
	return strings.TrimSuffix(strings.TrimPrefix(quotedPath, `"`), `"`)
}
