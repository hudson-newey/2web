package site

var paths = []string{}

func GetSitePaths() []string {
	return paths
}

func registerSitePage(path string) {
	paths = append(paths, path)
}
