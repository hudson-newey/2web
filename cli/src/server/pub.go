package server

import "fmt"

func Run(inPath string, outPath string) {
	fmt.Println(
		"Running in-built dev server.\n" +
			"To use Vite (for larger projects), run '2web template vite'.\n",
	)

	runDevServer(inPath, outPath)
}
