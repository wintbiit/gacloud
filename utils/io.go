package utils

import (
	"io"
)

func ToReadCloser(reader io.Reader) io.ReadCloser {
	if rc, ok := reader.(io.ReadCloser); ok {
		return rc
	}
	// return ioutil.NopCloser(reader)
	return io.NopCloser(reader)
}

func ToWriteCloser(writer io.Writer) io.WriteCloser {
	if wc, ok := writer.(io.WriteCloser); ok {
		return wc
	}
	return nil
}
