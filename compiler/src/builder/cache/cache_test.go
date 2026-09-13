package cache

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"hudson-newey/2web/src/filesystem"
)

func writeFile(t *testing.T, path string, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("failed to create directory for %s: %v", path, err)
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write %s: %v", path, err)
	}

	// File reads are cached for the lifetime of a compilation. Each key
	// computation in these tests simulates a fresh build, so the read cache
	// must be dropped to see the new file content.
	filesystem.ClearReadCache()
}

func keyFor(t *testing.T, inputPath string) string {
	t.Helper()

	key, err := CacheKey(inputPath, "/out/page.html")
	if err != nil {
		t.Fatalf("CacheKey returned an error: %v", err)
	}

	return key
}

// Regression tests for pages being served from cache after one of their
// compile time dependencies changed. The cache key is a fingerprint of the
// page's source and every file its compiled output depends on, so changing any
// of them must change the key.
func TestCacheKeyChangesWithPageContent(t *testing.T) {
	tempDir := t.TempDir()
	pagePath := filepath.Join(tempDir, "page.html")

	writeFile(t, pagePath, "<h1>v1</h1>")
	before := keyFor(t, pagePath)

	writeFile(t, pagePath, "<h1>v2</h1>")
	after := keyFor(t, pagePath)

	if before == after {
		t.Error("changing the page content did not change the cache key")
	}
}

func TestCacheKeyChangesWithLayoutChange(t *testing.T) {
	tempDir := t.TempDir()
	pagePath := filepath.Join(tempDir, "page.html")

	writeFile(t, pagePath, "<h1>Page</h1>")
	writeFile(t, filepath.Join(tempDir, "__layout.html"), "<html><slot></slot></html>")
	before := keyFor(t, pagePath)

	writeFile(t, filepath.Join(tempDir, "__layout.html"), "<html><body><slot></slot></body></html>")
	after := keyFor(t, pagePath)

	if before == after {
		t.Error("changing the shared __layout.html did not change the cache key")
	}
}

func TestCacheKeyChangesWithStyleSidecarChange(t *testing.T) {
	tempDir := t.TempDir()
	pagePath := filepath.Join(tempDir, "page.html")

	writeFile(t, pagePath, "<h1>Page</h1>")
	writeFile(t, filepath.Join(tempDir, "__style.css"), "h1 { color: red; }")
	before := keyFor(t, pagePath)

	writeFile(t, filepath.Join(tempDir, "__style.css"), "h1 { color: blue; }")
	after := keyFor(t, pagePath)

	if before == after {
		t.Error("changing the __style.css sidecar did not change the cache key")
	}
}

func TestCacheKeyChangesWithScriptSidecarChange(t *testing.T) {
	tempDir := t.TempDir()
	pagePath := filepath.Join(tempDir, "page.html")

	writeFile(t, pagePath, "<h1>Page</h1>")
	writeFile(t, filepath.Join(tempDir, "__script.ts"), "console.log('v1');")
	before := keyFor(t, pagePath)

	writeFile(t, filepath.Join(tempDir, "__script.ts"), "console.log('v2');")
	after := keyFor(t, pagePath)

	if before == after {
		t.Error("changing the __script.ts sidecar did not change the cache key")
	}
}

func TestCacheKeyChangesWithComponentImportChange(t *testing.T) {
	tempDir := t.TempDir()
	pagePath := filepath.Join(tempDir, "page.html")

	writeFile(t, pagePath, "<script compiled>\nimport Badge from \"components/badge.component.html\";\n</script>\n\n<Badge />")
	writeFile(t, filepath.Join(tempDir, "components", "badge.component.html"), "<span>v1</span>")
	before := keyFor(t, pagePath)

	writeFile(t, filepath.Join(tempDir, "components", "badge.component.html"), "<span>v2</span>")
	after := keyFor(t, pagePath)

	if before == after {
		t.Error("changing an imported component did not change the cache key")
	}
}

// Components can import other components, so a change to a transitive
// dependency must invalidate the page too.
func TestCacheKeyChangesWithTransitiveComponentChange(t *testing.T) {
	tempDir := t.TempDir()
	pagePath := filepath.Join(tempDir, "page.html")

	writeFile(t, pagePath, "<script compiled>\nimport Header from \"components/header.component.html\";\n</script>\n\n<Header />")
	writeFile(t, filepath.Join(tempDir, "components", "header.component.html"), "<script compiled>\nimport Title from \"./title.component.html\";\n</script>\n<Title />")
	writeFile(t, filepath.Join(tempDir, "components", "title.component.html"), "<title>v1</title>")
	before := keyFor(t, pagePath)

	writeFile(t, filepath.Join(tempDir, "components", "title.component.html"), "<title>v2</title>")
	after := keyFor(t, pagePath)

	if before == after {
		t.Error("changing a transitively imported component did not change the cache key")
	}
}

// ESM imports inside plain <script> tags are bundled by esbuild at compile
// time, so they are compile time dependencies as well.
func TestCacheKeyChangesWithEsmImportChange(t *testing.T) {
	tempDir := t.TempDir()
	pagePath := filepath.Join(tempDir, "page.html")

	writeFile(t, pagePath, "<script>\nimport { hello } from \"./scripts/esm.ts\";\n</script>")
	writeFile(t, filepath.Join(tempDir, "scripts", "esm.ts"), "export const hello = 'v1';")
	before := keyFor(t, pagePath)

	writeFile(t, filepath.Join(tempDir, "scripts", "esm.ts"), "export const hello = 'v2';")
	after := keyFor(t, pagePath)

	if before == after {
		t.Error("changing an ESM import did not change the cache key")
	}
}

// Dependencies that the compiler doesn't read must not invalidate the page.
func TestCacheKeyIgnoresUnrelatedFiles(t *testing.T) {
	tempDir := t.TempDir()
	pagePath := filepath.Join(tempDir, "page.html")

	writeFile(t, pagePath, "<h1>Page</h1>")
	writeFile(t, filepath.Join(tempDir, "__layout.html"), "<slot></slot>")
	before := keyFor(t, pagePath)

	writeFile(t, filepath.Join(tempDir, "unrelated.html"), "<p>not a dependency</p>")
	after := keyFor(t, pagePath)

	if before != after {
		t.Error("changing an unrelated file changed the cache key")
	}
}

func TestCacheKeyIsDeterministic(t *testing.T) {
	tempDir := t.TempDir()
	pagePath := filepath.Join(tempDir, "page.html")

	writeFile(t, pagePath, "<h1>Page</h1>")
	writeFile(t, filepath.Join(tempDir, "__layout.html"), "<slot></slot>")
	writeFile(t, filepath.Join(tempDir, "components", "badge.component.html"), "<span>v1</span>")

	first := keyFor(t, pagePath)
	second := keyFor(t, pagePath)

	if first != second {
		t.Error("CacheKey is not deterministic for identical inputs")
	}
}

func TestCacheKeyIncludesOutputPath(t *testing.T) {
	tempDir := t.TempDir()
	pagePath := filepath.Join(tempDir, "page.html")

	writeFile(t, pagePath, "<h1>Page</h1>")

	first, err := CacheKey(pagePath, "/out/one.html")
	if err != nil {
		t.Fatalf("CacheKey returned an error: %v", err)
	}

	second, err := CacheKey(pagePath, "/out/two.html")
	if err != nil {
		t.Fatalf("CacheKey returned an error: %v", err)
	}

	if first == second {
		t.Error("cache keys for different output paths should differ")
	}

	if !strings.HasPrefix(first, pagePath+"->"+"\""+"/out/one.html") && !strings.Contains(first, "/out/one.html") {
		t.Errorf("cache key %q does not reference its output path", first)
	}
}

// A page that can't be read can't be fingerprinted, so it must be treated as
// uncached (and uncacheable) instead of panicking or serving stale output.
func TestCacheKeyFailsForMissingPage(t *testing.T) {
	missingPage := filepath.Join(t.TempDir(), "does-not-exist.html")

	if _, err := CacheKey(missingPage, "/out/page.html"); err == nil {
		t.Error("CacheKey should return an error for a missing page file")
	}
}
