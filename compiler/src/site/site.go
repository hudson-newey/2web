package site

import "hudson-newey/2web/src/models"

var paths = []models.SitePath{}

func GetSitePaths() []models.SitePath {
	return paths
}

func registerSitePage(path models.SitePath) {
	paths = append(paths, path)
}
