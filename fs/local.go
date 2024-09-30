package fs

import (
	"context"
	"github.com/wintbiit/gacloud/utils"
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

func (l *localFileProvider) Get(ctx context.Context, fileSum string) (io.ReadSeekCloser, bool, error) {
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

func (l *localFileProvider) GetRanged(ctx context.Context, fileSum string, start, end int64) (io.ReadCloser, bool, error) {
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

	if _, err = f.Seek(start, 0); err != nil {
		return nil, false, err
	}

	limitReader := io.LimitReader(f, end-start)
	return utils.WithCloser(limitReader, f.Close), true, nil
}

func (l *localFileProvider) Put(ctx context.Context, reader io.Reader) (string, error) {
	tmpFile, err := os.CreateTemp(l.MountDir, "tmp")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if _, err = io.Copy(tmpFile, reader); err != nil {

		return "", err
	}

	if _, err = tmpFile.Seek(0, 0); err != nil {
		return "", err
	}

	fileSum := utils.Md5Sum(tmpFile)
	if err = os.MkdirAll(path.Join(l.MountDir, fileSum[:2], fileSum[2:4]), 0755); err != nil {
		return "", err
	}

	if err = tmpFile.Close(); err != nil {
		return "", err
	}

	if err = os.Rename(tmpFile.Name(), path.Join(l.MountDir, fileSum[:2], fileSum[2:4], fileSum)); err != nil {
		return "", err
	}

	return fileSum, nil
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
