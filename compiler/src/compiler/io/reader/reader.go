package reader

import (
	"bufio"
	"strings"
)

type Reader struct {
	FilePath string
	Reader   *bufio.Reader
}

func NewReader(filePath string, content string) *Reader {
	stringReader := strings.NewReader(content)
	reader := bufio.NewReader(stringReader)

	return &Reader{
		FilePath: filePath,
		Reader:   reader,
	}
}
