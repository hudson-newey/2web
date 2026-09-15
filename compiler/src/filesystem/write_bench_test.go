package filesystem

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func BenchmarkWriteOldDirect(b *testing.B) {
	dir := b.TempDir()
	content := bytes.Repeat([]byte("a"), 5*1024)
	path := fmt.Sprintf("%s/file-%d.txt", dir, os.Getpid())
	_ = os.WriteFile(path, content, 0644)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		const openMode = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
		file, err := os.OpenFile(path, openMode, 0644)
		if err != nil {
			b.Fatal(err)
		}
		if _, err := file.Write(content); err != nil {
			b.Fatal(err)
		}
		file.Close()
	}
}

func BenchmarkWriteNewAtomic(b *testing.B) {
	dir := b.TempDir()
	content := bytes.Repeat([]byte("a"), 5*1024)
	path := fmt.Sprintf("%s/file-%d.txt", dir, os.Getpid())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := writeFileAtomic(fileJob{Content: content, OutputPath: path}); err != nil {
			b.Fatal(err)
		}
	}
}
