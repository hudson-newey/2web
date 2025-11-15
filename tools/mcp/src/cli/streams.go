package cli

import "io"

func readStream(stream io.ReadCloser) string {
	outputBytes := make([]byte, 0)
	buf := make([]byte, 1024)
	for {
		n, err := stream.Read(buf)
		if n > 0 {
			outputBytes = append(outputBytes, buf[:n]...)
		}
		if err != nil {
			break
		}
	}

	return string(outputBytes)
}
