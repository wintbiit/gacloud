package fs

import (
	"context"
	"io"
)

type FileProvider interface {
	Get(ctx context.Context, fileSum string) (io.Reader, bool, error)
	Put(ctx context.Context, fileSum string) (io.Writer, error)
	Delete(ctx context.Context, fileSum string) error
	Exists(ctx context.Context, fileSum string) (bool, error)
}

type FileProviderFactory func([]byte) (FileProvider, error)

var fileProviderFactories = make(map[string]FileProviderFactory)
var fileProviderConfigs = make(map[string]interface{})

func RegisterFileProvider(name string, factory FileProviderFactory, config interface{}) {
	fileProviderFactories[name] = factory
	fileProviderConfigs[name] = config
}
