package fs

import (
	"context"
	"io"
)

type FileProvider interface {
	Get(ctx context.Context, fileSum string) (io.ReadSeekCloser, bool, error)
	Put(ctx context.Context, reader io.Reader) (string, error)
	Delete(ctx context.Context, fileSum string) error
	Exists(ctx context.Context, fileSum string) (bool, error)
}

type RangeFileProvider interface {
	GetRanged(ctx context.Context, fileSum string, start, end int64) (io.ReadCloser, bool, error)
}

type FileProviderFactory func([]byte) (FileProvider, error)

var (
	fileProviderFactories = make(map[string]FileProviderFactory)
	fileProviderConfigs   = make(map[string]interface{})
)

func RegisterFileProvider(name string, factory FileProviderFactory, config interface{}) {
	fileProviderFactories[name] = factory
	fileProviderConfigs[name] = config
}

func GetFileProviderFactory(name string) FileProviderFactory {
	return fileProviderFactories[name]
}

func ListFileProviderTypes() []string {
	var names []string
	for name := range fileProviderFactories {
		names = append(names, name)
	}
	return names
}
