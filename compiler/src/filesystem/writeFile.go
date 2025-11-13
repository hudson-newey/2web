package filesystem

import (
	"os"
	"path/filepath"
	"syscall"
	"time"
)

// Overwrites or creates a file at the given path with the given content and
// creates any necessary directories.
func WriteFile(content []byte, outputPath string) {
	if err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm); err != nil {
		panic(err)
	}

	// We use a lot of un-awaited goroutines here because writing to the files
	// should not be a blocking operation that stops the compiler from proceeding.
	writeNonBlocking(content, outputPath)
}

func CreateFile(outputPath string) {
	if err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm); err != nil {
		panic(err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
}

// A fast non-blocking file write operation that does NOT spawn a goroutine.
// This means that we can write to files without blocking the main thread or
// using a thread that would better be used for more expensive operations.
// TODO: Replace AI code
func writeNonBlocking(content []byte, outputPath string) {
	// If the file exists already, we delete it first.
	// TODO: Remove this hack. This was needed so that the old document is not
	// contaminated with the new document.
	syscall.Unlink(outputPath)

	fd, err := syscall.Open(outputPath, syscall.O_WRONLY|syscall.O_CREAT|syscall.O_NONBLOCK, 0644)
	if err != nil {
		panic(err)
	}
	defer syscall.Close(fd)

	totalWritten := 0

	for totalWritten < len(content) {
		n, err := syscall.Write(fd, content[totalWritten:])
		if err != nil {
			if err == syscall.EAGAIN || err == syscall.EWOULDBLOCK {
				// Resource temporarily unavailable, try again after a pause
				time.Sleep(100 * time.Millisecond)
				continue
			} else {
				panic(err)
			}
		}
		totalWritten += n
	}
}
