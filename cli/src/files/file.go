package files

type File struct {
	Path        string
	Content     string
	IsDirectory bool
	Children    []File
}
