package server

import (
	"context"

	"github.com/wintbiit/gacloud/fs"
	"github.com/wintbiit/gacloud/model"
	"github.com/wintbiit/gacloud/utils"
)

const DefaultFileProviderId = 0

func (s *GaCloudServer) RefreshFileProviders(ctx context.Context) error {
	providers, err := setupFileProviders(ctx, s.db)
	if err != nil {
		return err
	}

	s.fileProviders = providers
	return nil
}

// GetDefaultProviderID returns default file provider id for user
// priority: user provider > group provider > default provider(0)
func (s *GaCloudServer) GetDefaultProviderID(ctx context.Context, u *model.User) uint {
	// if defined user provider, return it
	var uf model.UserFileProvider
	err := s.db.WithContext(ctx).Where("user_id = ?", u.ID).First(&uf).Error
	if err == nil {
		return uf.FileProviderID
	}

	// if defined group provider, return it
	var gf model.GroupFileProvider
	err = s.db.WithContext(ctx).
		Joins("JOIN user_groups ON user_groups.group_id = group_file_providers.group_id").
		Where("user_groups.user_id = ?", u.ID).First(&gf).Error
	if err == nil {
		return gf.FileProviderID
	}

	// return default provider
	return DefaultFileProviderId
}

func (s *GaCloudServer) GetProvider(providerId uint) (fs.FileProvider, error) {
	p, ok := s.fileProviders[providerId]
	if !ok {
		return nil, utils.ErrorFileProviderNotFound
	}

	return p, nil
}

func (s *GaCloudServer) SetGroupFileProvider(ctx context.Context, groupId uint, providerId uint) error {
	gf := &model.GroupFileProvider{
		GroupID:        groupId,
		FileProviderID: providerId,
	}

	err := s.db.WithContext(ctx).Create(gf).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *GaCloudServer) SetUserFileProvider(ctx context.Context, userId uint, providerId uint) error {
	uf := &model.UserFileProvider{
		UserID:         userId,
		FileProviderID: providerId,
	}

	err := s.db.WithContext(ctx).Create(uf).Error
	if err != nil {
		return err
	}

	return nil
}
