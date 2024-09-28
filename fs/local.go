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
	RegisterFileProvider(localFileProviderName, NewLocalFileProvider, LocalFileProviderConfig{})
}

type LocalFileProviderConfig struct {
	MountDir string `json:"mount_dir"`
}

type localFileProvider struct {
	LocalFileProviderConfig
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

func (l *localFileProvider) Put(ctx context.Context, fileSum string, reader io.Reader) error {
	if len(fileSum) < 4 {
		return nil
	}

	p := path.Join(l.MountDir, fileSum[:2], fileSum[2:4], fileSum)
	if err := os.MkdirAll(path.Dir(p), 0o755); err != nil {
		return err
	}

	// if file exists, not dir and size > 0, regard as already exists, return nil
	stat, err := os.Stat(p)
	if err == nil && !stat.IsDir() && stat.Size() > 0 {
		return nil
	}

	if err == nil && stat.IsDir() {
		return os.ErrExist
	}

	f, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = io.Copy(f, reader); err != nil {
		return err
	}

	return nil
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

func NewLocalFileProvider(config []byte) (FileProvider, error) {
	var cfg LocalFileProviderConfig
	if err := json.Unmarshal(config, &cfg); err != nil {
		return nil, err
	}

	return NewLocalFileProviderWithConfig(cfg), nil
}

func NewLocalFileProviderWithConfig(cfg LocalFileProviderConfig) FileProvider {
	return &localFileProvider{cfg}
}
