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

type readerWithCloser struct {
	io.Reader
	Closer func() error
}

func (r *readerWithCloser) Close() error {
	if r.Closer == nil {
		return nil
	}
	return r.Closer()
}

func WithCloser(reader io.Reader, closer func() error) io.ReadCloser {
	return &readerWithCloser{
		Reader: reader,
		Closer: closer,
	}
}
