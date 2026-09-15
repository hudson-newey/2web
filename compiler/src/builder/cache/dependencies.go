package cache

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"hudson-newey/2web/src/cli"
	"hudson-newey/2web/src/filesystem"
)

// Files that the compiler implicitly pulls into a page when they sit next to
// it in the same directory.
//
// These are resolved relative to the page's own directory, mirroring:
//   - __layout.html:  1-preprocessor/layouts.go (expandLayout)
//   - __style.css:    5-templating/assets.go (addRouteAssets)
//   - __script.ts/js: 5-templating/assets.go (addRouteAssets)
var sidecarDependencies = []string{
	"__layout.html",
	"__style.css",
	"__script.ts",
	"__script.js",
}

// Matches static ESM style imports in page and dependency sources.
//
// This covers both component imports inside <script compiled> blocks (which
// are inlined at compile time by the import node) and plain ESM/TypeScript
// imports inside <script> blocks (which are bundled at compile time by
// esbuild). Both are compile time dependencies, so a change to the imported
// file must invalidate the page's cache entry.
var importPattern = regexp.MustCompile(`(?s)import\b[^;]*?from\s*(?:"([^"]+)"|'([^']+)')`)

// File extensions that are scanned for transitive imports.
//
// Any dependency is hashed regardless of extension, but only text-like source
// files are searched for further imports. This keeps the dependency walk away
// from binary assets (e.g. imported pdfs) that can never contain imports.
var importableExtensions = map[string]bool{
	".html": true,
	".htm":  true,
	".2web": true,
	".md":   true,
	".ts":   true,
	".tsx":  true,
	".js":   true,
	".mjs":  true,
	".jsx":  true,
	".css":  true,
}

// The maximum import chain depth that the dependency walk will follow.
// This only exists as a safety net; a visited set already prevents cycles.
const maxImportDepth = 16

// discoverDependencies returns every file that a page's compiled output
// depends on, other than the page's own source file.
//
// Dependencies are discovered in two ways:
//
//  1. Convention: the layout and style/script sidecar files that sit next to
//     the page.
//  2. Import scanning: static imports in the page's source, resolved and
//     followed transitively (components can import other components and ESM
//     modules can import other ESM modules).
//
// The returned paths are sorted so that the cache fingerprint is stable
// regardless of discovery order.
func discoverDependencies(inputPath string) []string {
	visited := map[string]bool{inputPath: true}
	dependencies := []string{}

	// Layout and sidecar dependencies are pulled into the page by convention
	// (they sit next to it in the same directory), so they can't be discovered
	// by scanning the page's source for imports.
	pageDirectory := filepath.Dir(inputPath)
	for _, sidecarName := range sidecarDependencies {
		candidate := filepath.Join(pageDirectory, sidecarName)
		if candidate == inputPath || visited[candidate] {
			continue
		}

		if _, err := os.Stat(candidate); err != nil {
			continue
		}

		visited[candidate] = true
		dependencies = append(dependencies, candidate)
	}

	discoverFrom(inputPath, &dependencies, visited, 0)
	sort.Strings(dependencies)

	return dependencies
}

func discoverFrom(filePath string, dependencies *[]string, visited map[string]bool, depth int) {
	if depth >= maxImportDepth {
		return
	}

	content, err := filesystem.ReadFile(filePath)
	if err != nil {
		return
	}

	for _, importPath := range findImports(string(content)) {
		for _, candidate := range resolveImport(filePath, importPath) {
			if visited[candidate] {
				continue
			}

			if _, err := os.Stat(candidate); err != nil {
				// The import target doesn't exist. The compiler will surface
				// this as a compiler error during compilation, so there is no
				// dependency to track here.
				visited[candidate] = true
				continue
			}

			visited[candidate] = true
			*dependencies = append(*dependencies, candidate)

			if canContainImports(candidate) {
				discoverFrom(candidate, dependencies, visited, depth+1)
			}
		}
	}
}

func findImports(content string) []string {
	matches := importPattern.FindAllStringSubmatch(content, -1)

	imports := []string{}
	for _, match := range matches {
		if len(match) < 3 {
			continue
		}

		// The import path can be single or double quoted.
		importPath := match[1]
		if importPath == "" {
			importPath = match[2]
		}

		if importPath != "" {
			imports = append(imports, importPath)
		}
	}

	return imports
}

// resolveImport returns the possible on-disk locations of an import path,
// in the same resolution order that the compiler uses at compile time.
//
// Component imports (see 4-parser/nodes/import.node.go) are resolved relative
// to the importing file's directory. ESM imports are bundled by esbuild, which
// resolves relative imports against the compiler's input path, so the input
// root is offered as a fallback.
func resolveImport(sourcePath string, importPath string) []string {
	if filepath.IsAbs(importPath) {
		return []string{filepath.Clean(importPath)}
	}

	sourceDirectory := filepath.Dir(sourcePath)
	candidates := []string{filepath.Join(sourceDirectory, importPath)}

	if inputRoot := inputRootDirectory(); inputRoot != "" && inputRoot != sourceDirectory {
		candidates = append(candidates, filepath.Join(inputRoot, importPath))
	}

	seen := map[string]bool{}
	resolved := []string{}
	for _, candidate := range candidates {
		candidate = filepath.Clean(candidate)
		if !seen[candidate] {
			seen[candidate] = true
			resolved = append(resolved, candidate)
		}
	}

	return resolved
}

// inputRootDirectory returns the directory that relative ESM imports are
// resolved against (the compiler's input path, or its parent when a single
// file is compiled).
func inputRootDirectory() string {
	inputPath := cli.GetArgs().InputPath
	if inputPath == "" {
		return ""
	}

	root := inputPath
	if info, err := os.Stat(root); err == nil && !info.IsDir() {
		root = filepath.Dir(root)
	}

	return root
}

func canContainImports(filePath string) bool {
	return importableExtensions[strings.ToLower(filepath.Ext(filePath))]
}
