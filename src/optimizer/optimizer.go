package optimizer

func OptimizeContent(content string) string {
	minifiedHtml := minifyHtml(content)
	return minifiedHtml
}
