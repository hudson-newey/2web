package filesystem

import (
	"os"
	"path/filepath"
	"syscall"
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

// A fast non-blocking file write operation that does NOT spawn a goroutine.
// This means that we can write to files without blocking the main thread or
// using a thread that would better be used for more expensive operations.
func writeNonBlocking(content []byte, outputPath string) {
	const fileMode int =
	// Open in write-only mode
	syscall.O_WRONLY |
		// Create the file if it does not exist
		syscall.O_CREAT |
		// Truncate the file if it already exists. This ensures that old content
		// is removed even if the new content is smaller.
		syscall.O_TRUNC |
		// We open the file in a non-blocking mode to prevent blocking the main
		// thread.
		// Warning: This does not guarantee that the write operation is completed
		// if another process is holding a lock on the file.
		// However, it does give use faster writes in most cases for a very small
		// risk of incomplete writes in rare cases that can usually be explained
		// by external user-level factors.
		syscall.O_NONBLOCK

	file, err := syscall.Open(outputPath, fileMode, 0644)
	if err != nil {
		panic(err)
	}
	defer syscall.Close(file)

	totalWritten := 0
	for totalWritten < len(content) {
		n, err := syscall.Write(file, content[totalWritten:])
		if err != nil {
			panic(err)
		}

		totalWritten += n
	}
}
