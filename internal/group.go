package internal

import (
	"context"
	"github.com/wintbiit/gacloud/model"
)

func (s *GaCloudServer) IsUserInGroup(ctx context.Context, user *model.User, groupId int64) bool {
	return false
}

func (s *GaCloudServer) UserAddGroup(ctx context.Context, user *model.User, group *model.Group) {
}

func (s *GaCloudServer) UserRemoveGroup(ctx context.Context, user *model.User, group *model.Group) {
}
