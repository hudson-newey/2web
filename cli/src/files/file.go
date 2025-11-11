package files

// TODO: Refactor this data model
type File struct {
	Path string

	// If this file is not a directory, you must set either Content or CopyBinary.
	// If this is a directory, you must set the Children property.
	// Note the Children property is not needed for files.
	IsDirectory bool
	Children    []File

	// These fields are mutually exclusive.
	Content      string
	CopyFromPath string
}
