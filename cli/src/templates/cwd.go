package templates

import (
	"log"
	"os"
)

func getCwd() string {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return path
}
