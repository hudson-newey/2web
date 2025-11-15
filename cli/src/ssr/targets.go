package ssr

import "os"

const SsrTargetDir string = "./server/"

func HasSsrTarget() bool {
	_, err := os.Stat(SsrTargetDir + "ssr.ts")
	return err == nil
}
