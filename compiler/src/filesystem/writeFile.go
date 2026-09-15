package filesystem

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"

	"hudson-newey/2web/src/cli"
)

type fileJob struct {
	Content    []byte
	OutputPath string
}

var fileQueue = make(chan fileJob, 128)

// pendingWrites tracks the number of file writes that have been queued to the
// file writer thread, but have not yet been flushed to disk.
//
// Because the file writer is asynchronous (see WriteFile), the compiler would
// otherwise exit before the writer thread has finished draining the write
// queue. This would silently drop queued writes or leave partially written
// files on disk. WaitFileWriter is used to drain the queue before the compiler
// reports a successful build.
var pendingWrites sync.WaitGroup

var failedWritesMutex sync.Mutex
var failedWrites []string

// TempFileSuffix is appended to output paths while a file is being written.
// Once the file has been fully written and closed, the temp file is renamed to
// the output path.
//
// Writing files through a temp file ensures that an output path never contains
// a partially written file. This is important because the build cache uses the
// existence of an output file to determine if an asset needs recompiling.
const TempFileSuffix = ".2webc-tmp"

// Overwrites or creates a file at the given path with the given content and
// creates any necessary directories.
//
// Unless the --serial flag is set, writes are handed off to a dedicated file
// writer thread so that the compiler is not blocked on disk I/O. Because that
// write is asynchronous, callers must call WaitFileWriter before claiming that
// the build output is complete (e.g. before the compiler exits).
func WriteFile(content []byte, outputPath string) {
	// Because this WriteFile function should be used by all parts of the
	// compiler, I can suppress any writes here when the --dry-run flag is set.
	if cli.GetArgs().DryRun {
		return
	}

	job := fileJob{
		content,
		outputPath,
	}

	if cli.GetArgs().Serial {
		if err := writeFileAtomic(job); err != nil {
			recordFailedWrite(job.OutputPath, err)
		}
		return
	}

	pendingWrites.Add(1)
	fileQueue <- job
}

// WaitFileWriter blocks until every file write queued through WriteFile has
// been flushed to disk.
//
// It is a no-op when the compiler is running in serial mode (writes are already
// synchronous) and in --dry-run mode (no writes are ever queued).
func WaitFileWriter() {
	pendingWrites.Wait()
}

// FailedWrites returns the output paths of all writes that failed since the
// compiler started. The returned slice is a copy and is safe to read while the
// file writer thread is still running.
func FailedWrites() []string {
	failedWritesMutex.Lock()
	defer failedWritesMutex.Unlock()

	paths := make([]string, len(failedWrites))
	copy(paths, failedWrites)
	return paths
}

func InitFileWriter() {
	if !cli.GetArgs().Serial {
		go fileWriterWorker()
	}
}

func fileWriterWorker() {
	// Since we'll be doing a lot of writes, keep this thread alive for the full
	// duration of the application so that we don't have to keep creating /
	// destroying the file writer thread.
	// TODO: Can we use a message-based system instead of spin locking here?
	for {
		for job := range fileQueue {
			if err := writeFileAtomic(job); err != nil {
				recordFailedWrite(job.OutputPath, err)
			}
			pendingWrites.Done()
		}
	}
}

// writeFileAtomic writes the given job to disk without ever exposing a
// partially written file at the job's output path.
//
// The content is written to a temp file in the same directory and is then
// renamed over the output path. Because the rename is only performed after the
// temp file has been fully written and closed, an interrupted build can never
// leave a truncated file at the output path.
// createdDirectories caches the directories that have already been created.
//
// Every file write used to call MkdirAll, which stats every path component of
// the directory on every write. With hundreds of output files that is
// thousands of redundant stat chains per build.
var createdDirectories sync.Map

func writeFileAtomic(job fileJob) error {
	outputDirectory := filepath.Dir(job.OutputPath)

	if _, created := createdDirectories.Load(outputDirectory); !created {
		const dirMode os.FileMode = os.ModeDir | 0755
		if err := os.MkdirAll(outputDirectory, dirMode); err != nil {
			return err
		}

		createdDirectories.Store(outputDirectory, struct{}{})
	}

	tempPath := job.OutputPath + TempFileSuffix

	// Open in write-only mode, creating the file if it does not exist and
	// truncating it if it does. This ensures that old content is removed even
	// if the new content is smaller.
	file, err := os.OpenFile(tempPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	if err = writeAll(file, job.Content); err != nil {
		file.Close()
		os.Remove(tempPath)
		return err
	}

	// Close returns an error if the file's buffered data could not be flushed
	// to the disk. If we ignored this error, the rename below could promote an
	// incomplete file to the output path.
	if err = file.Close(); err != nil {
		os.Remove(tempPath)
		return err
	}

	if err = os.Rename(tempPath, job.OutputPath); err != nil {
		os.Remove(tempPath)
		return err
	}

	return nil
}

// writeAll writes the entire content buffer to the given writer.
//
// io.Writer implementations (including os.File) are allowed to perform short
// writes, so a single Write call is not guaranteed to write the entire buffer.
func writeAll(writer io.Writer, content []byte) error {
	for len(content) > 0 {
		n, err := writer.Write(content)
		if n > 0 {
			content = content[n:]
		}

		if err != nil {
			return err
		}

		if n == 0 {
			return io.ErrShortWrite
		}
	}

	return nil
}

func recordFailedWrite(outputPath string, err error) {
	log.Printf("Error writing file %s: %v", outputPath, err)

	failedWritesMutex.Lock()
	defer failedWritesMutex.Unlock()
	failedWrites = append(failedWrites, outputPath)
}

// WriteFileSync writes a file synchronously and atomically (through a temp
// file that is renamed over the output path).
//
// Use this for files whose presence must be guaranteed as soon as the
// function returns (e.g. the server entry point that the cli executes after
// the build), and for reporting write failures.
func WriteFileSync(content []byte, outputPath string) error {
	if cli.GetArgs().DryRun {
		return nil
	}

	if cli.GetArgs().Serial {
		return writeFileAtomic(fileJob{Content: content, OutputPath: outputPath})
	}

	// Synchronous writes bypass the write queue: the file must exist when this
	// function returns.
	if err := os.MkdirAll(filepath.Dir(outputPath), os.ModeDir|0755); err != nil {
		return err
	}

	tempPath := outputPath + TempFileSuffix

	file, err := os.OpenFile(tempPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	if err = writeAll(file, content); err != nil {
		file.Close()
		os.Remove(tempPath)
		return err
	}

	if err = file.Close(); err != nil {
		os.Remove(tempPath)
		return err
	}

	return os.Rename(tempPath, outputPath)
}
