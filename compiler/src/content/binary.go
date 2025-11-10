package content

type BinaryFile interface {
	Data() []byte

	// This is a function in case you want to add some hashing or path information
	// e.g. for storing the resolution of an image in the path.
	OutputPath() string
}
