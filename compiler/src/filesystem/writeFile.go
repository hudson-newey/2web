package filesystem

import (
	"hudson-newey/2web/src/cli"
	"log"
	"os"
	"path/filepath"
	"syscall"
)

type fileJob struct {
	Content    []byte
	OutputPath string
}

var fileQueue = make(chan fileJob, 128)

// Overwrites or creates a file at the given path with the given content and
// creates any necessary directories.
// There is always one file writer thread which handles writing to files.
// By using a writer thread, we don't block the compiler on writing files.
func WriteFile(content []byte, outputPath string) {
	// Because this WriteFile function should be used by all parts of the
	// compiler, I can suppress any writes here when the --dry-run flag is set.
	if cli.GetArgs().DryRun {
		return
	}

	fileQueue <- fileJob{
		content,
		outputPath,
	}
}

func InitFileWriter() {
	go fileWriterWorker()
}

func fileWriterWorker() {
	for job := range fileQueue {
		const dirMode os.FileMode = os.ModeDir | 0755
		if err := os.MkdirAll(filepath.Dir(job.OutputPath), dirMode); err != nil {
			log.Printf("Error creating directory for %s: %v", job.OutputPath, err)
			continue
		}

		const openMode int =
		// Open in write-only mode
		syscall.O_WRONLY |
			// Create the file if it does not exist
			syscall.O_CREAT |
			// Truncate the file if it already exists. This ensures that old content
			// is removed even if the new content is smaller.
			syscall.O_TRUNC

		file, err := os.OpenFile(job.OutputPath, openMode, 0644)
		if err != nil {
			log.Printf("Error opening file %s: %v", job.OutputPath, err)
			continue
		}

		if _, err = file.Write(job.Content); err != nil {
			log.Printf("Error writing to file %s: %v", job.OutputPath, err)
		}
		defer file.Close()
	}
}
