package server

import (
	"context"
	"io"
	"sync"

	"github.com/wintbiit/gacloud/fs"
	"github.com/wintbiit/gacloud/model"
	"github.com/wintbiit/gacloud/utils"
)

var (
	fileProviders map[int64]fs.FileProvider
	filePool      = sync.Pool{
		New: func() interface{} {
			return &model.File{}
		},
	}
	listFilePool = sync.Pool{
		New: func() interface{} {
			return &model.ListFile{}
		},
	}
)

func (s *GaCloudServer) GetFileByPath(ctx context.Context, p string) (*model.File, func(), error) {
	file := filePool.Get().(*model.File)
	clean := func() {
		filePool.Put(file)
	}

	err := s.db.WithContext(ctx).Where("path = ?", p).First(file).Error
	if err != nil {
		return nil, clean, err
	}

	return file, clean, nil
}

func (s *GaCloudServer) GetFileById(ctx context.Context, id int64) (*model.File, func(), error) {
	file := filePool.Get().(*model.File)
	clean := func() {
		filePool.Put(file)
	}

	err := s.db.WithContext(ctx).First(file, id).Error
	if err != nil {
		return nil, clean, err
	}

	return file, clean, nil
}

func (s *GaCloudServer) GetFileBySum(ctx context.Context, sum string) (*model.File, func(), error) {
	file := filePool.Get().(*model.File)
	clean := func() {
		filePool.Put(file)
	}

	err := s.db.WithContext(ctx).Where("sum = ?", sum).First(file).Error
	if err != nil {
		return nil, clean, err
	}

	return file, clean, nil
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

func (s *GaCloudServer) ListFiles(ctx context.Context, operator *model.User, dir string) ([]*model.ListFile, func(), error) {
	dir = utils.CleanDirPath(dir)

	tx := s.db.WithContext(ctx).Table("files AS f").
		Joins("LEFT JOIN user_groups AS ug ON files.owner_id = user_groups.group_id AND files.owner_type = ?", model.FileOwnerTypeGroup).
		Where("(f.owner_type = ? AND f.owner_id = ?) OR (f.owner_type = ? AND ug.user_id = ?)", model.FileOwnerTypeUser, operator.ID, model.FileOwnerTypeGroup, operator.ID).
		Group("f.path").
		Select("f.path, f.size, f.mime, f.owner_type, f.owner_id").
		Order("f.path")

	files := make([]*model.ListFile, tx.RowsAffected)
	for i := range files {
		files[i] = listFilePool.Get().(*model.ListFile)
	}

	clean := func() {
		for _, f := range files {
			listFilePool.Put(f)
		}
	}

	err := tx.Scan(&files).Error
	if err != nil {
		return nil, nil, err
	}

	return files, clean, nil
}
