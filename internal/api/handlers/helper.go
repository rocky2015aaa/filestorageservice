package handlers

import (
	"bytes"
	"io"
)

// splitFile splits the file into `n` parts
func splitFile(file io.Reader, n int) ([][]byte, error) {
	var parts [][]byte
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, file)
	if err != nil {
		return nil, err
	}

	fileSize := buf.Len()
	partSize := fileSize / n
	for i := 0; i < n; i++ {
		start := i * partSize
		end := start + partSize
		if i == n-1 {
			end = fileSize
		}
		parts = append(parts, buf.Bytes()[start:end])
	}
	return parts, nil
}
