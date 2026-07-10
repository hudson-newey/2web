package builder

import (
	"hudson-newey/2web/src/cli"
	"hudson-newey/2web/src/parallel/threadPool"
	"hudson-newey/2web/src/routing"
	"math"
	"runtime"
)

// Inspired by Astro 7 Queued rendering
// https://astro.build/blog/astro-7/#:~:text=Queued%20Rendering,-Queued%20rendering
//
// Instead of recursively compiling web pages, we maintain a queue of pages that
// we want to compile and a thread pool that reads the compile queue.
// This prevents us from repeatedly spinning up new threads and we can use flat
// itteration opposed to recursion.

var buildPool *threadPool.ThreadPool

func startBuildPool() {
	const assignmentRatio float64 = 0.8 // 80%
	threadSize := int(math.Floor(float64(runtime.NumCPU()) * assignmentRatio))
	threadSize = max(threadSize, 1)
	if cli.GetArgs().Serial {
		threadSize = 1
	}

	buildPool = threadPool.NewThreadPool(threadSize)
}

func AddCompilationStep(filePath string) {
	args := cli.GetArgs()
	buildPool.Enqueue(func() {
		if routing.IsLayoutFile(filePath) {
			return
		}

		compileAndWritePage(
			filePath,
			outputFileName(args.InputPath, args.OutputPath, filePath),
		)
	})
}
