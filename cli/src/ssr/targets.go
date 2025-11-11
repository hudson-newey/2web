package ssr

import "os"

func HasSsrTarget() bool {
	_, err := os.Stat("./server/ssr.ts")
	return err == nil
}
