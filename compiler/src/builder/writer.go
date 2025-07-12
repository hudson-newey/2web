package builder

func compileAndWriteFile(inputPath string, outputPath string) {
	compiledPage := buildToPage(inputPath)
	compiledPage.Write(outputPath)
}
