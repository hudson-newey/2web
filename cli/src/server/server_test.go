package server

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDiffSnapshotReportsCreatedAndModifiedFiles(t *testing.T) {
	root := t.TempDir()

	write := func(name string, content string) {
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	write("a.css", "v1")
	write("__2web.debug.json", "{}")
	before := snapshotDirectory(root)

	// Rewriting a file with the same content (the compiler does this for the
	// sitemap and other bookkeeping files) must not count as a change.
	write("a.css", "v1")
	write("__2web.debug.json", fmt.Sprintf("{\"built\": \"%d\"}", time.Now().UnixNano()))
	write("b.js", "new")
	write("nested/c.css", "new")

	after := snapshotDirectory(root)
	changed := diffSnapshot(before, after)

	if len(changed) != 2 {
		t.Fatalf("expected 2 changed files (the rewritten bookkeeping file and the untouched file are excluded), got %v", changed)
	}
}

func TestSnapshotDirectoryHandlesMissingRoot(t *testing.T) {
	snapshot := snapshotDirectory(filepath.Join(t.TempDir(), "does-not-exist"))

	if len(snapshot) != 0 {
		t.Errorf("expected an empty snapshot for a missing root, got %v", snapshot)
	}
}
