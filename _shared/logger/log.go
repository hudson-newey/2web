package logger

import "fmt"

func Println(msg string) (n int, err error) {
	return fmt.Println(msg)
}
