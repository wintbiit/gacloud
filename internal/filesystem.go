package internal

import (
	"context"
	"github.com/wintbiit/gacloud/fs"
	"github.com/wintbiit/gacloud/model"
	"github.com/wintbiit/gacloud/utils"
	"io"
)

var fileProviders map[int64]fs.FileProvider

func (s *GaCloudServer) GetFileByPath(ctx context.Context, p string) *model.File {
	return nil
}

func (s *GaCloudServer) GetFileById(ctx context.Context, id int64) *model.File {
	return nil
}

func (s *GaCloudServer) GetFileBySum(ctx context.Context, sum string) *model.File {
	return nil
}

func (s *GaCloudServer) GetFileReader(ctx context.Context, f *model.File) (io.Reader, bool, error) {
	provider, ok := fileProviders[f.ProviderId]
	if !ok {
		return nil, false, utils.ErrorFileProviderNotFound
	}

	return provider.Get(ctx, f.Sum)
}

func (s *GaCloudServer) GetFileWriter(ctx context.Context, f *model.File) (io.Writer, error) {
	provider, ok := fileProviders[f.ProviderId]
	if !ok {
		return nil, utils.ErrorFileProviderNotFound
	}

	return provider.Put(ctx, f.Sum)
}

func (s *GaCloudServer) DeleteFile(ctx context.Context, f *model.File) error {
	provider, ok := fileProviders[f.ProviderId]
	if !ok {
		return utils.ErrorFileProviderNotFound
	}

	return provider.Delete(ctx, f.Sum)
}

func (s *GaCloudServer) FileExists(ctx context.Context, f *model.File) (bool, error) {
	provider, ok := fileProviders[f.ProviderId]
	if !ok {
		return false, utils.ErrorFileProviderNotFound
	}

	return provider.Exists(ctx, f.Sum)
}
