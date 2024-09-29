package server

import (
	"context"

	"github.com/wintbiit/gacloud/model"
)

func (s *GaCloudServer) IsUserInGroup(ctx context.Context, user *model.User, groupId uint) bool {
	count := s.db.WithContext(ctx).Model(&model.UserGroup{}).Where("user_id = ? AND group_id = ?", user.ID, groupId).RowsAffected
	return count > 0
}

func (s *GaCloudServer) UserAddGroup(ctx context.Context, user *model.User, group *model.Group) error {
	ug := &model.UserGroup{
		User:  *user,
		Group: *group,
	}

	if err := s.db.WithContext(ctx).Create(ug).Error; err != nil {
		return err
	}

	return nil
}

func (s *GaCloudServer) UserRemoveGroup(ctx context.Context, user *model.User, group *model.Group) error {
	if err := s.db.WithContext(ctx).Where("user_id = ? AND group_id = ?", user.ID, group.ID).Delete(&model.UserGroup{}).Error; err != nil {
		return err
	}

	return nil
}
