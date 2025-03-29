package compiler

import (
	"log"
)

func Compile(path string, content string) string {
	log.Println("\t-", path)
	return content
}
