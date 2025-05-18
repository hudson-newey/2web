package packages

func HasPackageJson() bool {
	return fileExists("package.json")
}
