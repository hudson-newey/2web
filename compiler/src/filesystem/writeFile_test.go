package filesystem

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

// Regression test for builds exiting before the asynchronous file writer
// finished draining its queue. The compiler used to silently drop queued
// writes, leaving assets missing from the build output.
func TestWaitFileWriterFlushesAllQueuedWrites(t *testing.T) {
	startFileWriterThread()
	tempDir := t.TempDir()

	const writeCount = 500
	contents := make(map[string][]byte)

	for i := range writeCount {
		outputPath := filepath.Join(tempDir, fmt.Sprintf("asset-%03d.txt", i))
		content := []byte(fmt.Sprintf("content of asset %d", i))
		contents[outputPath] = content

		WriteFile(content, outputPath)
	}

	WaitFileWriter()

	for outputPath, expectedContent := range contents {
		written, err := os.ReadFile(outputPath)
		if err != nil {
			t.Fatalf("queued write to %s was never flushed to disk: %v", outputPath, err)
		}

		if !bytes.Equal(written, expectedContent) {
			t.Errorf("file %s contains %q, expected %q", outputPath, written, expectedContent)
		}
	}

	leftoverTempFiles, err := filepath.Glob(filepath.Join(tempDir, "*"+TempFileSuffix))
	if err != nil {
		t.Fatalf("failed to glob for leftover temp files: %v", err)
	}

	if len(leftoverTempFiles) != 0 {
		t.Errorf("found %d leftover temp files: %v", len(leftoverTempFiles), leftoverTempFiles)
	}
}

// The write queue must be safe when many compilation workers enqueue writes at
// the same time.
func TestConcurrentWritersDoNotCorruptWrites(t *testing.T) {
	startFileWriterThread()
	tempDir := t.TempDir()

	const writeCount = 500
	var wg sync.WaitGroup

	for i := range writeCount {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			outputPath := filepath.Join(tempDir, fmt.Sprintf("concurrent-%03d.txt", i))
			WriteFile([]byte(fmt.Sprintf("content of asset %d", i)), outputPath)
		}(i)
	}

	wg.Wait()
	WaitFileWriter()

	for i := range writeCount {
		outputPath := filepath.Join(tempDir, fmt.Sprintf("concurrent-%03d.txt", i))
		expectedContent := fmt.Sprintf("content of asset %d", i)

		written, err := os.ReadFile(outputPath)
		if err != nil {
			t.Fatalf("write to %s was never flushed to disk: %v", outputPath, err)
		}

		if string(written) != expectedContent {
			t.Errorf("file %s contains %q, expected %q", outputPath, written, expectedContent)
		}
	}
}

// startFileWriterThread ensures exactly one file writer worker is draining the
// write queue for the whole test binary.
//
// The production build pipeline starts the writer thread through
// InitFileWriter, but tests enqueue writes directly.
var writerThreadOnce sync.Once

func startFileWriterThread() {
	writerThreadOnce.Do(InitFileWriter)
}

// A new, smaller write must fully replace the previous content of a file. The
// old implementation opened the output with O_TRUNC but ignored short writes,
// which could leave stale bytes at the end of the file.
func TestWriteFileReplacesFileContent(t *testing.T) {
	startFileWriterThread()
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "style.css")

	WriteFile([]byte("body { background: red; }\n/* a lot of stale content that should be truncated */"), outputPath)
	WaitFileWriter()

	replacement := []byte("body { background: blue; }\n")
	WriteFile(replacement, outputPath)
	WaitFileWriter()

	written, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read %s: %v", outputPath, err)
	}

	if !bytes.Equal(written, replacement) {
		t.Errorf("file %s contains %q, expected %q", outputPath, written, replacement)
	}
}

// Atomic writes must never expose a partially written file at the output path.
// The build cache treats the existence of an output file as "this asset is up
// to date", so a truncated file would stick around until the source changes.
func TestFailedWriteDoesNotTruncateExistingOutput(t *testing.T) {
	startFileWriterThread()
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "index.html")

	originalContent := []byte("<html></html>\n")
	if err := os.WriteFile(outputPath, originalContent, 0644); err != nil {
		t.Fatalf("failed to seed output file: %v", err)
	}

	// Force a write failure by pointing the write at a directory conflict
	// (the parent path of the output is a file, so MkdirAll fails).
	conflictingPath := filepath.Join(tempDir, "blocker")
	if err := os.WriteFile(conflictingPath, []byte("not a directory"), 0644); err != nil {
		t.Fatalf("failed to create blocking file: %v", err)
	}

	before := FailedWrites()
	if err := writeFileAtomic(fileJob{Content: []byte("should never be written"), OutputPath: filepath.Join(conflictingPath, "nested", "out.html")}); err == nil {
		t.Fatal("expected write into a file-as-directory path to fail")
	}

	// The queued write path must also record the failure.
	WriteFile([]byte("should never be written"), filepath.Join(conflictingPath, "nested", "queued.html"))
	WaitFileWriter()

	written, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read %s: %v", outputPath, err)
	}

	if !bytes.Equal(written, originalContent) {
		t.Errorf("failed write corrupted an existing output file: %q", written)
	}

	if len(FailedWrites()) != len(before)+1 {
		t.Errorf("failed write was not recorded in FailedWrites")
	}
}

// io.Writer implementations are allowed to perform short writes, so the write
// loop must retry until the entire buffer has been written.
type shortWriteWriter struct {
	buffer       bytes.Buffer
	maxWriteSize int
}

func (w *shortWriteWriter) Write(content []byte) (int, error) {
	if len(content) > w.maxWriteSize {
		content = content[:w.maxWriteSize]
	}

	return w.buffer.Write(content)
}

func TestWriteAllRetriesShortWrites(t *testing.T) {
	writer := &shortWriteWriter{maxWriteSize: 7}
	content := []byte(strings.Repeat("abcdefgh", 32))

	if err := writeAll(writer, content); err != nil {
		t.Fatalf("writeAll returned an error: %v", err)
	}

	if !bytes.Equal(writer.buffer.Bytes(), content) {
		t.Error("writeAll did not write the entire buffer")
	}
}
