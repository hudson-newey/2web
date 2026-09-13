package preprocessor

import (
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
func expandComponents(filePath string, content string, visiting map[string]bool) string {
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

		// Only markup files are components. e.g. ESM imports of .ts files are
		// bundled by esbuild and must be left alone.
		if !assets.IsMarkupFile(componentPath) {
			continue
		}

		// Components are inlined by replacing their selector
		// (e.g. <Header />). A component without a selector in the importing
		// file isn't used, so inlining it would only pollute the page with its
		// script variables.
		selector := "<" + importName + " />"
		if !strings.Contains(content, selector) {
			continue
		}

		visiting[componentPath] = true
		expandedComponent := expandComponents(componentPath, string(componentContent), visiting)
		delete(visiting, componentPath)

		content = strings.Replace(content, statement, "", 1)
		content = strings.ReplaceAll(content, selector, expandedComponent)
	}

	return content
}

func unquoteImportPath(quotedPath string) string {
	return strings.TrimSuffix(strings.TrimPrefix(quotedPath, `"`), `"`)
}
