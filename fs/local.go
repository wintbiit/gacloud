package fs

import (
	"context"
	"io"
	"os"
	"path"

	"github.com/goccy/go-json"
)

const localFileProviderName = "local"

func init() {
	RegisterFileProvider(localFileProviderName, newLocalFileProvider, localFileProviderConfig{})
}

type localFileProviderConfig struct {
	MountDir string `json:"mount_dir"`
}

type localFileProvider struct {
	localFileProviderConfig
}

func (l *localFileProvider) Get(ctx context.Context, fileSum string) (io.ReadCloser, bool, error) {
	if len(fileSum) < 4 {
		return nil, false, nil
	}

	p := path.Join(l.MountDir, fileSum[:2], fileSum[2:4], fileSum)
	f, err := os.OpenFile(p, os.O_RDONLY, 0)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false, nil
		}
		return nil, false, err
	}

	return f, true, nil
}

func (l *localFileProvider) Put(ctx context.Context, fileSum string) (io.WriteCloser, error) {
	if len(fileSum) < 4 {
		return nil, nil
	}

	p := path.Join(l.MountDir, fileSum[:2], fileSum[2:4], fileSum)
	if err := os.MkdirAll(path.Dir(p), 0o755); err != nil {
		return nil, err
	}

	f, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (l *localFileProvider) Delete(ctx context.Context, fileSum string) error {
	if len(fileSum) < 4 {
		return nil
	}

	p := path.Join(l.MountDir, fileSum[:2], fileSum[2:4], fileSum)
	if err := os.Remove(p); err != nil {
		return err
	}

	return nil
}

func (l *localFileProvider) Exists(ctx context.Context, fileSum string) (bool, error) {
	if len(fileSum) < 4 {
		return false, nil
	}

	p := path.Join(l.MountDir, fileSum[:2], fileSum[2:4], fileSum)
	if _, err := os.Stat(p); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func newLocalFileProvider(config []byte) (FileProvider, error) {
	var cfg localFileProviderConfig
	if err := json.Unmarshal(config, &cfg); err != nil {
		return nil, err
	}

	return &localFileProvider{cfg}, nil
}
