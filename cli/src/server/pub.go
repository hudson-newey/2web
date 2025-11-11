package server

import "fmt"

func Run(inPath string, outPath string) {
	// We use a double newline here to separate the initial message from
	// subsequent log output.
	fmt.Printf(
		"Running in-built dev server.\n" +
			"To use Vite (for larger projects), run '2web template vite'.\n\n",
	)

	runDevServer(inPath, outPath)
}
