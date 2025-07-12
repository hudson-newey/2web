package preprocessor

type isStable = bool

func ProcessStaticSite(filePath string, content string) (string, isStable) {
	ssgResult := content

	return ssgResult, true
}
